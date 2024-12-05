面试题目：

Airbnb数据采集器

你的任务是创建一个工具，用于在Airbnb上收集特定类别的所有预订信息。请使用生产者-消费者模型来设计这个工具。

可交付成果:

您需要交付一个名为airbnb-cli的可执行文件，该文件可以通过命令行运行。

airbnb-cli的功能：

启动一个使用者服务器，它可以配置并发程度和侦听端口。

    例子:
**airbnb-cli start consumer --workers=10 --queue=<your-queue-server>**

启动一个生产者向消费者发送任务。

    例子:
**airbnb-cli start producer --data tasks.json --queue=<your-queue-server>**

要求:

找到一个允许生产者和消费者交换信息的队列解决方案。

允许多个消费者和生产者同时运行。

确保每个任务只执行一次，没有重复。

使用MySQL存储从web上抓取的结果。自己设计数据库模式。

示例tasks.json:

```Plain Text
{
  "tasks": [
    {
      "name": "Europe travel",
      "url": "https://www.airbnb.com/s/Europe/homes?tab_id=home_tab&refinement_paths%5B%5D=%2Fhomes&flexible_trip_lengths%5B%5D=one_week&monthly_start_date=2024-12-01&monthly_length=3&monthly_end_date=2025-03-01&price_filter_input_type=0&channel=EXPLORE&place_id=ChIJhdqtz4aI7UYRefD8s-aZ73I&date_picker_type=calendar&source=structured_search_input_header&search_type=filter_change",
      "headers": {
        "user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X x.y; rv:42.0) Gecko/20100101 Firefox/42.0"
      }
    }
  ]
}

class="gsgwcjk" div
遍历内部div 获取class="itemListElement" 内部 meta标签itemprop="url" 进入详情页
aria-label="下一个" a标签进入下一页

```


- name：任务名称。

- url: Airbnb上的目录页。

- headers：用于访问Airbnb的头文件配置。

任务详细信息:

对于每个任务，您需要抓取指定目录下的所有子页面并保存其预订信息。

您必须将以下信息存储到数据库中

1. 酒店的名字【data-plugin-in-point-id="OVERVIEW_DEFAULT_V2" div下class=hpipapi 的h2标签】

2. 多少星【客房推荐的 data-plugin-in-point-id="GUEST_FAVORITE_BANNER" div 下的 class=a8jhwcl div 下的 aria-hidden="true" div内容 或者 非推荐的直接使用class=r1lutz1s div的内容】

3. 价格【class=_11jcbg2 的span内容】

4. 税前价格【class=_1avmy66 div下的 class=_j1kt73的span内容】

5. 入住日期【data-testid="change-dates-checkIn" div内容】

6. 退房日期【data-testid="change-dates-checkOut" div内容】

7. 客人【class= _7pspom div下的 class=_j1kt73 的span内容】

提交说明

- 存储库:

    - 提供一个包含完整代码库的GitHub存储库链接。

- README.md:

    - 包括文档部分中概述的所有必要部分。

- 代码质量:

    - 按照最佳实践确保代码干净、易读和可维护。

    - 在整个应用程序中进行正确的错误处理和输入验证。

- 测试:

    - 验证所有端点工作正常。

    - 使用提供的示例API请求演示工作流。

    - 显示由中间件生成的示例日志。

    
```Plain Text
提取信息 class="gsgwcjk"的div内部的div遍历内部div 获取class="itemListElement" 内部 meta标签itemprop="url" 进入详情页，点击aria-label="下一个"的a标签进入下一页，直到没有aria-label="下一个"的a标签完成所有分页的遍历
详情页数据获取的规则根据:
1. 酒店的名字【data-plugin-in-point-id="OVERVIEW_DEFAULT_V2" div下class=hpipapi 的h2标签】
2. 多少星【客房推荐的 data-plugin-in-point-id="GUEST_FAVORITE_BANNER" div 下的 class=a8jhwcl div 下的 aria-hidden="true" div内容 或者 非推荐的直接使用class=r1lutz1s div的内容】
3. 价格【class=_11jcbg2 的span内容】
4. 税前价格【class=_1avmy66 div下的 class=_j1kt73的span内容】
5. 入住日期【data-testid="change-dates-checkIn" div内容】
6. 退房日期【data-testid="change-dates-checkOut" div内容】
7. 客人【class= _7pspom div下的 class=_j1kt73 的span内容】
```

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