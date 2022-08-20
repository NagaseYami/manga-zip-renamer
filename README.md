# manga-zip-renamer

这是一个当我用`Calibre`+`DeDRM`+`KindleUnpack`破解了Kindle电子书并提取了zip文件之后  
需要批量重命名zip文件和zip里的图片文件时所使用的小工具  
※破解方案具体请看[这里](http://blog.nagaseyami.com/posts/Kindle%E7%94%B5%E5%AD%90%E4%B9%A6%E5%9B%BE%E6%BA%90%E6%8F%90%E5%8F%96/)  

## 大致内容

1. 递归寻找给定文件夹下所有`zip`文件
1. 在zip文件同目录里寻找`metadata.opf`文件并从中获取Kindle书本文件名和作者名
1. 用简陋的方式尝试从Kindle书本文件名中提取出书本名称和卷标
1. 最后以`[作者名] 书本名 第N卷`的方式重命名
1. 把原zip内的所有文件以`0001`~`9999`的方式重命名后复制进新的zip文件
1. 输出在`被D&D文件夹的同目录下/output/[作者名] 书本名`文件夹中
