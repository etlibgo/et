package et

//邮箱规则
const REGEX_EMAIL_RULE = "^\\w+((-\\w+)|(\\.\\w+))*\\@[A-Za-z0-9]+((\\.|-)[A-Za-z0-9]+)*\\.[A-Za-z0-9]+$"
//手机号规则
const REGEX_PHONE_RULE = "^1[3456789]\\d{9}$"
//手机号规则（进阶）
const REGEX_MOBILE_PHONE_RULE = "^((13[0-9])|(14[5|7])|(15([0-3]|[5-9]))|(18[0,1,2,3,5-9])|(17[0-9]))\\d{8}$"
//多手机号规则,逗号相隔
const REGEX_MULTI_PHONE_RULE = "^1[34578][0-9]{9}(,1[34578][0-9]{9})*$"
//多手机号规则,逗号相隔
const REGEX_QQ_RULE = "[1-9][0-9]{4,14}"
//包含三位小数
const REGEX_THREE_NUMERIC_RULE = "^(([1-9]{1}\\d*)|([0]{1}))(\\.(\\d){0,3})?$"
//验证是数字和字母含大小写
const REGEX_NUM_OR_LETTER_RULE = "^[0-9a-zA-Z]*$"
//纯数字验证
const REGEX_NUM_RULE = "^[0-9]*$"
//固话验证
const REGEX_TEL_RULE = "^[1-9]{1}[0-9]{5,8}$"
//固话验证带区号
const REGEX_TEL_ZONE_RULE = "^[0][1-9][0-9]{1,2}-[0-9]{5,10}$"
//中文验证
const REGEX_CHINESE_RULE = "^[\\u4e00-\\u9fa5]*$"
//密码验证，6-20位、由数字与字母构成、至少包含一位数字和一位字母、不可带空格
const REGEX_PASSWORD_RULE = "^[0-9a-zA-Z](?=.*\\d.*)(?=.*[a-zA-Z].*)(?!.*\\s.*).[A-Za-z0-9]{6,20}$"
