/*
 * @Description:放置配置config相关的信息对象
 */
package config

//
// 放置config配置相关的信息对象
//

type Config struct {
	Server Server `json:"server"`
}

//
// 服务器相关的配置信息
//
type Server struct {
	Port int16 `json:"port"`
}
