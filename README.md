Go Newrelic Package
====================

Provides:
- Wrapper for Newrelic SDK 
- Library for recording transactions
- Library for pushing events to Newrelic Insights 



Install
-------

```
go get github.com/barnjamin/newrelic
```

cd $GOHOME/github.com/barnjamin/newrelic

Download newrelic sdk by following directions here:
https://docs.newrelic.com/docs/agents/agent-sdk/installation-configuration/installing-agent-sdk

If you untar the contents to this directory under `newrelic_sdk` the C -L flag will not need to be updated in wrapper.go, otherwise change that path to whatever the new path should be


Update the `NEWRELIC_LICENSE_KEY` in wrapper.go to use transactions and set the Go Version

Update the `NEWRELIC_LICENSE_KEY` in events.go to use the event tracker and set the url for the insight collector


```
cd examples
go run transactionExample.go
```

If that complains about not being able to find a shared lib make sure your LD_LIBRARY_PATH is set to the same directory as your -L flag in wrapper.go and that the shared libraries are actually there. 



