## README

csv_loader 是一个将csv表格标准化为接口给go调用的建议工具

下面给出几个参数类型:

+ -c, --csvPath: 默认不需要配置, 只有当embed开启时候自动配置值为module, 为嵌入资源的二级目录文件名
+ -d, --definedFilePath: 默认不需要配置, 值为"02_csv_defined.go", 为csvloader类型声明文件名
+ -e, --embed: 默认不需要配置, 值为false, 当开启时候改为csv表嵌入二进制文件, 注意开启后无法热更新
+ -f, --filePath: 默认不需要配置, 值为"01_csv_table.go", 为csvloader类型
+ -F, --formulaFilePath: 默认不需要配置, 值为"03_formula.go", 为csvloader的公式导入
+ -g, --genPath: **核心配置参数, 没有默认值, 需要生成文件地址**
+ -h, --help                     help for csv_loader
+ -p, --packageName: 默认不需要配置, 值为"csv", 为package包名
+ -t, --templatePath: 默认不需要配置, 当配置模版地址的时候, 使用定制模版而不是嵌入二进制文件内的模版文件
+ -w, --workPath: **核心配置参数, 为csv文件地址**
+ -x, --xlsxFilePath: 默认不需要配置, 当开启的时候可以使用xlsx进行loader解析, 注意需要依赖excel2csv, 值为xlsx地址
