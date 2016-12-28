# go-timing

在日常的设计中我们会出数十万Go程，出现超时的问题，不能使用大量chan管道，所有借鉴其他作者siddontang的思路；
http://blog.csdn.net/siddontang/article/details/18370541?utm_source=tuicool&utm_medium=referral
在他的基础上，去除Mutex，只使用CAP操作提升性能。
