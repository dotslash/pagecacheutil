# pagecacheutil

A less featureful version of https://github.com/hoytech/vmtouch written in go as a learning excercise. It lets you 
1) see what parts of a file are in [page cache](https://en.wikipedia.org/wiki/Page_cache)
2) add a file into page cache
3) evict all aprts of a file from page cache (the system call for eviction is different in osx and linux. see oscompat directory)

### Screenshot with usage

<img src=https://storage.googleapis.com/tmp-uploads-2adf0f005374/665cfbcb-ee73-4ea3-9ed4-1a0d989983dc.png>
