# Testkube Testsuites Run

## **Synopsis**

Starts a new Test Suite based on TestSuite Custom Resource name and returns the results to the console.

```
testkube testsuites run <testSuiteName> [flags]
```

## **Options**

```
  -a, --download-artifacts     Downlaod artifacts automatically.
      --download-dir string    Download directory (default "artifacts").
  -h, --help                   Help for run.
  -n, --name string            Execution name. If empty, will be autogenerated.
  -p, --param stringToString   Execution envs passed to the executor (default []).
  -f, --watch                  Watch for changes after start.
```

## **Options Inherited from Parent Commands**

```
      --analytics-enabled    Enable analytics (default "true").
  -c, --client string        Client used for connecting to testkube API one of proxy|direct (default "proxy").
      --go-template string   When choosing output==go, pass golang template (default "{{ . | printf \"%+v\"  }}").
  -s, --namespace string     Kubernetes namespace (default "testkube").
  -o, --output string        Output type - raw, json or go  (default "raw").
  -v, --verbose              Show additional debug messages.
```

## **SEE ALSO**

* [Testkube Testsuites](testkube_testsuites.md)	 - Testsuites management commands.

