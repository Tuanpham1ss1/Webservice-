package infrastructure

import (
	"flag"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"os"
)

const (
	APPPORT    = "APP_PORT"
	DBHOST     = "DB_HOST"
	DBPORT     = "DB_PORT"
	DBNAME     = "DB_NAME"
	DBUSER     = "DB_USER"
	DBPASSWORD = "DB_PASSWORD"

	HTTPSWAGGER = "HTTP_SWAGGER"

	PRIVATE_PATH = "PRIVATE_PATH"
	PUPLIC_PATH  = "PUBLIC_PATH"

	REDIS_URL = "REDIS_URL"
)

var (
	dbName     string
	appPort    string
	dbHost     string
	dbPort     string
	dbUser     string
	dbPassword string

	httpSwagger string

	db *gorm.DB

	privatePath string
	publicPath  string

	redisURL    string
	redisClient *redis.Client

	encodeAuth *JWTAuth
	decodeAuth *JWTAuth

	ErrLog  *log.Logger
	InfoLog *log.Logger
)

func GetHTTPSwagger() string {
	return httpSwagger
}
func GetDBName() string {
	return dbName
}
func GetAppPort() string {
	return appPort
}
func GetDB() *gorm.DB {
	return db
}
func GetRedisClient() *redis.Client {
	return redisClient
}
func GetEncodeAuth() *JWTAuth {
	return encodeAuth
}
func GetDecodeAuth() *JWTAuth {
	return decodeAuth
}

// func GetEncodeAuth()
func gotDotEnvVariable(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		ErrLog.Println("Error loading .env file")
	}
	return os.Getenv(key)
}
func getStringEnvParameter(envParam string, defaultValue string) string {
	if value, ok := os.LookupEnv(envParam); ok {
		return value
	}
	return defaultValue
}
func loadEnvParameters(dbNameArg string, dbPwdArg string) {
	//root
	root, _ := os.Getwd()

	appPort = getStringEnvParameter(APPPORT, "8080")
	dbHost = getStringEnvParameter(DBHOST, "localhost")
	dbPort = getStringEnvParameter(DBPORT, "5432")
	dbName = getStringEnvParameter(DBNAME, dbNameArg)
	dbUser = getStringEnvParameter(DBUSER, "postgres")
	dbPassword = getStringEnvParameter(DBPASSWORD, dbPwdArg)
	httpSwagger = getStringEnvParameter(HTTPSWAGGER, gotDotEnvVariable(HTTPSWAGGER))

	privatePath = getStringEnvParameter(PRIVATE_PATH, root+"/infrastructure/private.pem")
	publicPath = getStringEnvParameter(PUPLIC_PATH, root+"/infrastructure/public.pem")

	redisURL = getStringEnvParameter(REDIS_URL, "localhost:6379")
}

func init() {
	ErrLog = log.New(os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLog = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)

	//nếu true sẽ tạo database, nếu false sẽ không tạo database
	var initDB bool
	flag.BoolVar(&initDB, "initDB", false, "allow recreate model database in postgres")

	var dbNameArg string
	flag.StringVar(&dbNameArg, "dbname", "webservice", "database name")

	var dbPwdArg string
	flag.StringVar(&dbPwdArg, "dbpwd", "12345", "database password")

	flag.Parse() // Phân tích các cờ từ dòng lệnh
	loadEnvParameters(dbNameArg, dbPwdArg)

	// Init database
	if err := InitDatabase(initDB); err != nil {
		log.Println("error initialize database: ", err)
		panic(err)
	}
	// Init redis
	if err := InitRedis(); err != nil {
		log.Println("error initialize redis: ", err)
		panic(err)
	}
	// Init JWT
	if err := loadAuthToken(); err != nil {
		log.Println("error load auth token: ", err)
		panic(err)
	}
}
