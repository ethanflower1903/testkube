package result

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
)

const CollectionName = "results"

func NewMongoRespository(db *mongo.Database) *MongoRepository {
	return &MongoRepository{
		Coll: db.Collection(CollectionName),
	}
}

type MongoRepository struct {
	Coll *mongo.Collection
}

func (r *MongoRepository) Get(ctx context.Context, id string) (result testkube.Execution, err error) {
	err = r.Coll.FindOne(ctx, bson.M{"id": id}).Decode(&result)
	return
}

func (r *MongoRepository) GetByNameAndTest(ctx context.Context, name, testName string) (result testkube.Execution, err error) {
	err = r.Coll.FindOne(ctx, bson.M{"name": name, "testname": testName}).Decode(&result)
	return
}

func (r *MongoRepository) GetNewestExecutions(ctx context.Context, limit int) (result []testkube.Execution, err error) {
	result = make([]testkube.Execution, 0)
	resultLimit := int64(limit)
	opts := &options.FindOptions{Limit: &resultLimit}
	opts.SetSort(bson.D{{Key: "_id", Value: -1}})
	cursor, err := r.Coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return result, err
	}
	cursor.All(ctx, &result)
	return
}

func (r *MongoRepository) GetExecutions(ctx context.Context, filter Filter) (result []testkube.Execution, err error) {
	result = make([]testkube.Execution, 0)
	query, opts := composeQueryAndOpts(filter)

	cursor, err := r.Coll.Find(ctx, query, opts)
	if err != nil {
		return
	}
	err = cursor.All(ctx, &result)

	return
}

func (r *MongoRepository) GetExecutionTotals(ctx context.Context, paging bool, filter ...Filter) (totals testkube.ExecutionsTotals, err error) {
	var result []struct {
		Status string `bson:"_id"`
		Count  int32  `bson:"count"`
	}

	query := bson.M{}
	if len(filter) > 0 {
		query, _ = composeQueryAndOpts(filter[0])
	}

	pipeline := []bson.D{{{"$match", query}}}
	if len(filter) > 0 {
		pipeline = append(pipeline, bson.D{{"$sort", bson.D{{"starttime", -1}}}})
		if paging {
			pipeline = append(pipeline, bson.D{{"$skip", int64(filter[0].Page() * filter[0].PageSize())}})
			pipeline = append(pipeline, bson.D{{"$limit", int64(filter[0].PageSize())}})
		}
	}

	pipeline = append(pipeline, bson.D{{"$group", bson.D{{"_id", "$executionresult.status"}, {"count", bson.D{{"$sum", 1}}}}}})
	cursor, err := r.Coll.Aggregate(ctx, pipeline)
	if err != nil {
		return totals, err
	}
	err = cursor.All(ctx, &result)
	if err != nil {
		return totals, err
	}

	var sum int32

	// TODO: statuses are messy e.g. success==passed error==failed
	for _, o := range result {
		sum += o.Count
		switch testkube.TestSuiteExecutionStatus(o.Status) {
		case testkube.QUEUED_TestSuiteExecutionStatus:
			totals.Queued = o.Count
		case testkube.PENDING_TestSuiteExecutionStatus:
			totals.Pending = o.Count
		case testkube.SUCCESS_TestSuiteExecutionStatus:
			totals.Passed = o.Count
		case testkube.ERROR__TestSuiteExecutionStatus:
			totals.Failed = o.Count
		}
	}
	totals.Results = sum

	return
}

func (r *MongoRepository) GetLabels(ctx context.Context) (labels map[string][]string, err error) {
	var result []struct {
		Labels bson.M `bson:"labels"`
	}

	cursor, err := r.Coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &result)
	if err != nil {
		return nil, err
	}

	labels = map[string][]string{}
	for _, r := range result {
		for key, value := range r.Labels {
			if values, ok := labels[key]; !ok {
				labels[key] = []string{fmt.Sprint(value)}
			} else {
				for _, v := range values {
					if v == value {
						continue
					}
				}
				labels[key] = append(labels[key], fmt.Sprint(value))
			}
		}
	}
	return labels, nil
}

func (r *MongoRepository) Insert(ctx context.Context, result testkube.Execution) (err error) {
	_, err = r.Coll.InsertOne(ctx, result)
	return
}

func (r *MongoRepository) Update(ctx context.Context, result testkube.Execution) (err error) {
	_, err = r.Coll.ReplaceOne(ctx, bson.M{"id": result.Id}, result)
	return
}

func (r *MongoRepository) UpdateResult(ctx context.Context, id string, result testkube.ExecutionResult) (err error) {
	_, err = r.Coll.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": bson.M{"executionresult": result}})
	return
}

// StartExecution updates execution start time
func (r *MongoRepository) StartExecution(ctx context.Context, id string, startTime time.Time) (err error) {
	_, err = r.Coll.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": bson.M{"starttime": startTime}})
	return
}

// EndExecution updates execution end time
func (r *MongoRepository) EndExecution(ctx context.Context, id string, endTime time.Time, duration time.Duration) (err error) {
	_, err = r.Coll.UpdateOne(ctx, bson.M{"id": id}, bson.M{"$set": bson.M{"endtime": endTime, "duration": duration.String()}})
	return
}

func composeQueryAndOpts(filter Filter) (bson.M, *options.FindOptions) {
	query := bson.M{}
	opts := options.Find()
	startTimeQuery := bson.M{}

	if filter.TextSearchDefined() {
		query["$or"] = bson.A{
			bson.M{"testname": bson.M{"$regex": primitive.Regex{Pattern: filter.TextSearch(), Options: "i"}}},
			bson.M{"name": bson.M{"$regex": primitive.Regex{Pattern: filter.TextSearch(), Options: "i"}}},
		}
	}

	if filter.TestNameDefined() {
		query["testname"] = filter.TestName()
	}

	if filter.StartDateDefined() {
		startTimeQuery["$gte"] = filter.StartDate()
	}

	if filter.EndDateDefined() {
		startTimeQuery["$lte"] = filter.EndDate()
	}

	if len(startTimeQuery) > 0 {
		query["starttime"] = startTimeQuery
	}

	if filter.StatusDefined() {
		query["executionresult.status"] = filter.Status()
	}

	if filter.Selector() != "" {
		items := strings.Split(filter.Selector(), ",")
		for _, item := range items {
			elements := strings.Split(item, "=")
			if len(elements) == 2 {
				query["labels."+elements[0]] = elements[1]
			} else if len(elements) == 1 {
				query["labels."+elements[0]] = bson.M{"$exists": true}
			}
		}
	}

	if filter.TypeDefined() {
		query["testtype"] = filter.Type()
	}

	opts.SetSkip(int64(filter.Page() * filter.PageSize()))
	opts.SetLimit(int64(filter.PageSize()))
	opts.SetSort(bson.D{{"starttime", -1}})

	return query, opts
}
