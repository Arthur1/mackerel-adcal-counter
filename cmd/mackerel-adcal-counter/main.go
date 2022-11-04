package main

import (
	"context"
	"fmt"
	"mackerel-adcal-counter/config"
	"mackerel-adcal-counter/counter/qiita"
	"mackerel-adcal-counter/writer"
	writer_mackerel "mackerel-adcal-counter/writer/mackerel"
	writer_stdout "mackerel-adcal-counter/writer/stdout"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/mackerelio/mackerel-client-go"
)

type Event struct{}

type Response struct{}

func Handler(ctx context.Context, evt Event) (Response, error) {
	conf, err := newConfigFromEnv(ctx)
	if err != nil {
		return Response{}, err
	}

	c := qiita.NewCounter(conf.Url)
	result, err := c.Count()
	if err != nil {
		return Response{}, err
	}

	var w writer.Writer
	if conf.MackerelApiKeyIsSet() {
		w, err = writer_mackerel.NewWriter(conf.MackerelApiKey)
		if err != nil {
			return Response{}, err
		}
	} else {
		w = writer_stdout.NewWriter()
	}
	w.WriteWithContext(ctx, conf.Service, []*mackerel.MetricValue{
		{
			Name:  fmt.Sprintf("%s.entries", conf.MetricName),
			Time:  time.Now().Unix(),
			Value: result.Entries,
		},
		{
			Name:  fmt.Sprintf("%s.participants", conf.MetricName),
			Time:  time.Now().Unix(),
			Value: result.Participants,
		},
		{
			Name:  fmt.Sprintf("%s.subscribers", conf.MetricName),
			Time:  time.Now().Unix(),
			Value: result.Subscribers,
		},
	})
	return Response{}, nil
}

func main() {
	lambda.Start(Handler)
}

func newConfigFromEnv(ctx context.Context) (*config.Config, error) {
	var apiKey string
	var err error
	ssmKey := os.Getenv("MACKEREL_API_KEY_SSM")
	if ssmKey != "" {
		apiKey, err = fetchMackerelApiKey(ctx, ssmKey)
		if err != nil {
			return nil, err
		}
	}

	return config.NewConfig(os.Getenv("URL"), os.Getenv("SERVICE"), os.Getenv("METRIC_NAME"), apiKey)
}

func fetchMackerelApiKey(ctx context.Context, ssmKey string) (string, error) {
	sess := session.Must(session.NewSession())
	sess.Config.Region = aws.String("ap-northeast-1")
	client := ssm.New(sess)
	res, err := client.GetParameterWithContext(ctx, &ssm.GetParameterInput{
		Name:           aws.String(ssmKey),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", err
	}

	return *res.Parameter.Value, nil
}
