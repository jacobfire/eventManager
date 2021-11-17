package apiserver

import (
	"calendar/configs"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	SECRETKEY string = "HELLOJWT"
)

func auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := strings.Split(r.Header.Get("Authorization"), "Bearer ")
		if len(authHeader) != 2 {
			fmt.Println("Unathorized access")
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte("Unathorized access"))
			if err != nil {
				log.Println(err)
			}
			return
		}

		jwtToken := authHeader[1]
		redisConfig := configs.NewConfig().RedisConfig
		redisClient := redis.NewClient(&redis.Options{
			Addr:     redisConfig.Host + redisConfig.Port,
			Password: redisConfig.Password,
			DB:       redisConfig.DB,
		})

		res, err := redisClient.Get(jwtToken).Result()

		if err == redis.Nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte("Unathorized access"))
			if err != nil {
				log.Println(err)
			}
			return
		}

		if res == "" {
			fmt.Println("Unathorized access")
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte("Unathorized access"))
			if err != nil {
				log.Println(err)
			}
			return
		}
		//ctx := context.WithValue(r.Context(), "authtoken", SECRETKEY)
		//next.ServeHTTP(w, r.WithContext(ctx))
		//fmt.Println(authHeader)

		token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(SECRETKEY), nil
		})

		type propeties string
		var props propeties= "props"
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			ctx := context.WithValue(r.Context(), props, claims)
			// Access context values in handlers like this
			// props, _ := r.Context().Value("props").(jwt.MapClaims)

			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			fmt.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte("Unauthorized"))
			if err != nil {
				log.Println(err)
			}
		}
	})
}

// GenerateJWT gives us necessary token for the session
func GenerateJWT(email, role string) (string, error) {
	var mySigningKey = []byte(SECRETKEY)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	tokenString, err := token.SignedString(mySigningKey)
	redisConfig := configs.NewConfig().RedisConfig
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Host + redisConfig.Port,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})

	err = redisClient.Set(tokenString, email, 84600000000).Err()
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return tokenString, nil
}

func Logout(token string) error {
	redisConfig := configs.NewConfig().RedisConfig
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Host + redisConfig.Port,
		Password: redisConfig.Password,
		DB:       redisConfig.DB,
	})

	if len(token) == 0 {
		return errors.New("empty token parameter")
	}
	value, err := redisClient.Get(token).Result()
	fmt.Println("TOKENINREDIS = ", value)

	if err == redis.Nil {
		return errors.New("not existing token")
	}

	if len(value) == 0 {
		return errors.New("unavailable user")
	}
	err = redisClient.Del(token).Err()
	if err != nil {
		return err
	}

	return nil
}
