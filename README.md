# mit6.824-2017-mapreduce
> MIT6.824 2017课程作业的lab1，实现mapreduce。框架代码来自 [git://g.csail.mit.edu/6.824-golabs-2017](git://g.csail.mit.edu/6.824-golabs-2017)

# 1 MapReduce基础概念
map 负责分发。每个map任务通常处理一个文件，有多少个输入文件就有多少个map任务。而输出则根据reduce的数目确定有多少个输出。注意，相同的key肯定是输出到同一个reduce文件中，通过哈希算法来确认一个key应该分发到哪个中间文件。targetReduct = ihash(key) % nReduce。
总的中间文件数为 nMap * nReduce

reduce负责对map输出的文件进行处理，每个reduce处理自己负责的那些中间文件。负责处理的中间文件名可以由map个数和reduce number来确认。

# 2 代码实现
原始代码有两个目录：main和mapreduce。其中main目录下只需要关注`wc.go，ii.go以及一系列的txt文件`。 mapreduce目录则需要修改`schedule.go`，`common_map.go`和`common_reduce.go`文件。
## Part 1， Part 2
这两个实验都比较简单，单机的。Part 1实现common_map.go和common_reduce.go中的通用函数，而Part 2也只是统计单词数目。Part 2这里要注意的是，要使用实验文档中说的分割方法`unicode.IsLetter()`去分割单词，不然测试会无法通过。

## Part 3 
采用的分布式执行任务。先是master启动rpc服务，并调用schedule函数执行Task。注意的是，Worker启动后需要调用master的RPC方法Register注册到master，而master发现有了新的worker，会通知等待channel的协程进行任务操作，同一个worker需要处理多个任务，但是不能同时给一个worker分配多个任务。

注意在run中，先做map任务，做完map再做reduce。所以在多个worker做map任务的时候，需要等所有的map任务完成再继续reduce任务。这里用了 `WaitGroup`。

另外，分配worker的算法要注意，是每次有新的可能会分配到，而老的worker如果执行完了一次任务，则也要放回channel中以重复使用。

## Part 4 
worker失败的情况下的实验，也比较简单，在每个任务执行时加个for循环，如果成功则退出，否则重新取worker执行任务。

## Part 5
Part 5需要完成一个倒排索引，完成ii.go即可，跟wc.go类似，主要是要去掉重复以及对倒排列表排序就可以通过测试了。
