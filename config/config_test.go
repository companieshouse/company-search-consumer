package config_test

import (
	"encoding/json"
	"os"
	"strconv"

	"github.com/companieshouse/company-search-consumer/config"
	"regexp"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// key constants
const (
	CONSUMERGROUPNAMECONST         = `CONSUMER_GROUP_NAME`
	STREAMCOMPANYPROFILETOPICCONST = `STREAM_COMPANY_PROFILE_TOPIC`
	SCHEMAREGISTRYURLCONST         = `SCHEMA_REGISTRY_URL`
	STREAMINGBROKERADDRCONST       = `KAFKA_STREAMING_BROKER_ADDR`
	STREAMINGZOOKEEPERURLCONST     = `KAFKA_ZOOKEEPER_STREAMING_ADDR`
	ZOOKEEPERCHROOTECONST          = `KAFKA_ZOOKEEPER_CHROOT`
	INITIALOFFSETCONST             = `INITIAL_OFFSET`
	LOGTOPICCONST                  = `LOG_TOPIC`
	ALPHABETICALUPSERTURLCONST     = `ALPHABETICAL_UPSERT_URL`
	ADVANCEDUPSERTURLCONST         = `ADVANCED_UPSERT_URL`
	RETRYTHROTTLERATECONST         = `RETRY_THROTTLE_RATE_SECONDS`
	MAXRETRYATTEMPTSCONST          = `MAXIMUM_RETRY_ATTEMPTS`
	CHSAPIKEYCONST                 = `CHS_INTERNAL_API_KEY`
)

// value constants
const (
	consumerGroupNameConst         = `consumer-group-name`
	streamCompanyProfileTopicConst = `stream-company-profile-topic`
	schemaRegistryUrlConst         = `schema-registry-url`
	streamingBrokerAddrConst       = `streaming-broker-addr`
	streamingZookeeperUrlConst     = `streaming-zookeeper-url`
	zookeeperChrootConst           = `zookeeper-chroot`
	initialOffsetConst             = 123
	logTopicConst                  = `log-topic`
	alphabeticalUpsertUrlConst     = `alphabetical-upsert-url`
	advancedUpsertUrlConst         = `advanced-upsert-url`
	retryThrottleRateConst         = 456
	maxRetryAttemptsConst          = 789
	chsApiKeyConst                 = `chs-api-key`
)

func TestConfig(t *testing.T) {
	t.Parallel()
	os.Clearenv()
	var (
		err           error
		configuration *config.Config
		envVars       = map[string]string{
			CONSUMERGROUPNAMECONST:         consumerGroupNameConst,
			STREAMCOMPANYPROFILETOPICCONST: streamCompanyProfileTopicConst,
			SCHEMAREGISTRYURLCONST:         schemaRegistryUrlConst,
			STREAMINGBROKERADDRCONST:       streamingBrokerAddrConst,
			STREAMINGZOOKEEPERURLCONST:     streamingZookeeperUrlConst,
			ZOOKEEPERCHROOTECONST:          zookeeperChrootConst,
			INITIALOFFSETCONST:             strconv.Itoa(initialOffsetConst),
			LOGTOPICCONST:                  logTopicConst,
			ALPHABETICALUPSERTURLCONST:     alphabeticalUpsertUrlConst,
			ADVANCEDUPSERTURLCONST:         advancedUpsertUrlConst,
			RETRYTHROTTLERATECONST:         strconv.Itoa(retryThrottleRateConst),
			MAXRETRYATTEMPTSCONST:          strconv.Itoa(maxRetryAttemptsConst),
			CHSAPIKEYCONST:                 chsApiKeyConst,
		}
		builtConfig = config.Config{
			ConsumerGroupName:         consumerGroupNameConst,
			StreamCompanyProfileTopic: streamCompanyProfileTopicConst,
			SchemaRegistryURL:         schemaRegistryUrlConst,
			StreamingBrokerAddr:       []string{streamingBrokerAddrConst},
			StreamingZookeeperURL:     streamingZookeeperUrlConst,
			ZookeeperChroot:           zookeeperChrootConst,
			InitialOffset:             initialOffsetConst,
			LogTopic:                  logTopicConst,
			AlphabeticalUpsertURL:     alphabeticalUpsertUrlConst,
			AdvancedUpsertURL:         advancedUpsertUrlConst,
			RetryThrottleRate:         retryThrottleRateConst,
			MaxRetryAttempts:          maxRetryAttemptsConst,
			CHSAPIKey:                 chsApiKeyConst,
		}
		consumerGroupNameRegex          = regexp.MustCompile(consumerGroupNameConst)
		streamCompanyProfileTopicRegex = regexp.MustCompile(streamCompanyProfileTopicConst)
		schemaRegistryUrlRegex          = regexp.MustCompile(schemaRegistryUrlConst)
		streamingBrokerAddrRegex        = regexp.MustCompile(streamingBrokerAddrConst)
		streamingZookeeperUrlRegex      = regexp.MustCompile(streamingZookeeperUrlConst)
		zookeeperChrootRegex            = regexp.MustCompile(zookeeperChrootConst)
		initialOffsetRegex              = regexp.MustCompile(strconv.Itoa(initialOffsetConst))
		logTopicRegex                   = regexp.MustCompile(logTopicConst)
		alphabeticalUpsertUrlRegex      = regexp.MustCompile(alphabeticalUpsertUrlConst)
		advancedUpsertUrlRegex          = regexp.MustCompile(advancedUpsertUrlConst)
		retryThrottleRateRegex          = regexp.MustCompile(strconv.Itoa(retryThrottleRateConst))
		maxRetryAttemptsRegex           = regexp.MustCompile(strconv.Itoa(maxRetryAttemptsConst))
		chsApiKeyRegex                  = regexp.MustCompile(chsApiKeyConst)
	)

	// set test env variables
	for varName, varValue := range envVars {
		os.Setenv(varName, varValue)
		defer os.Unsetenv(varName)
	}

	Convey("Given an environment with no environment variables set", t, func() {

		Convey("Then configuration should be nil", func() {
			So(configuration, ShouldBeNil)
		})

		Convey("When the config values are retrieved", func() {

			Convey("Then there should be no error returned, and values are as expected", func() {
				configuration = config.Get()

				So(err, ShouldBeNil)
				So(configuration, ShouldResemble, &builtConfig)
			})

			Convey("The generated JSON string from configuration should not contain sensitive data", func() {
				jsonByte, err := json.Marshal(builtConfig)

				So(err, ShouldBeNil)
				So(consumerGroupNameRegex.Match(jsonByte), ShouldEqual, true)
				So(streamCompanyProfileTopicRegex.Match(jsonByte), ShouldEqual, true)
				So(schemaRegistryUrlRegex.Match(jsonByte), ShouldEqual, true)
				So(streamingBrokerAddrRegex.Match(jsonByte), ShouldEqual, true)
				So(streamingZookeeperUrlRegex.Match(jsonByte), ShouldEqual, true)
				So(zookeeperChrootRegex.Match(jsonByte), ShouldEqual, true)
				So(initialOffsetRegex.Match(jsonByte), ShouldEqual, true)
				So(logTopicRegex.Match(jsonByte), ShouldEqual, true)
				So(alphabeticalUpsertUrlRegex.Match(jsonByte), ShouldEqual, true)
				So(advancedUpsertUrlRegex.Match(jsonByte), ShouldEqual, true)
				So(retryThrottleRateRegex.Match(jsonByte), ShouldEqual, true)
				So(maxRetryAttemptsRegex.Match(jsonByte), ShouldEqual, true)
				So(chsApiKeyRegex.Match(jsonByte), ShouldEqual, false)
			})
		})
	})
}
