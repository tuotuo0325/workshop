windows本地环境无法开启hyper-v，无法搭建docker服务，没有mysql、redis、kafaka环境，采集数据用到的消息队列与mysql存储数据全部使用文件来操作。

# 技术栈选择

## 消息队列【可以考虑通过增加参数来选择具体服务并进行扩展】
    - kafaka    【公司存在kafaka服务可选择】
    - redis list 【小公司也会有redis服务】
    - 文件【环境不足，选择文件，当前使用文件】

## 数据库存储
    - mysql
    - 文件

## 每个任务只执行一次
    - 使用redis布隆过滤器过滤
    - 当前使用方案：使用文件记录任务执行状态

# 生产者
airbnb-cli.exe producer --data task.json --queue data/queue.data

# 消费者
airbnb-cli.exe consumer --workers 10 --queue data/queue.data --storage data/hotels.json