package main

import (
	"flag"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

func forceNewDeployment(svc *ecs.ECS, cluster string, service string) error {
	log.Printf("starting deployment of '%s' in '%s' cluster", service, cluster)

	input := &ecs.UpdateServiceInput{
		Cluster:            aws.String(cluster),
		Service:            aws.String(service),
		ForceNewDeployment: aws.Bool(true),
	}

	_, err := svc.UpdateService(input)
	if err != nil {
		return err
	}

	return nil
}

func describeService(svc *ecs.ECS, cluster string, service string) (*ecs.DescribeServicesOutput, error) {
	input := &ecs.DescribeServicesInput{
		Cluster: aws.String(cluster),
		Services: []*string{
			aws.String(service),
		},
	}

	result, err := svc.DescribeServices(input)
	if err != nil {
		return &ecs.DescribeServicesOutput{}, err
	}

	return result, nil
}

func getDeploymentStatusTilComplete(svc *ecs.ECS, cluster string, service string, interval int) {
	// Loop status til all deployments show completed
	for {
		// Get status of the deployment
		results, err := describeService(svc, cluster, service)
		if err != nil {
			log.Fatal(err)
		}

		deploymentComplete := true
		for _, service := range results.Services {
			for _, deployment := range service.Deployments {
				deploymentStatus := "complete"
				if *deployment.DesiredCount != *deployment.RunningCount {
					deploymentComplete = false
					deploymentStatus = "running"
				}

				log.Printf("deployment status, updated at: %s (%s)", deployment.UpdatedAt.Format("Mon Jan _2 15:04:05 2006"), deploymentStatus)
				log.Printf("deployment id: %s", *deployment.Id)
				log.Printf(
					"desired count: %d, running count: %d, pending count: %d",
					*deployment.DesiredCount,
					*deployment.RunningCount,
					*deployment.PendingCount,
				)
			}
		}

		// Create a seperator in logging for better readability
		log.Println("-")

		if deploymentComplete {
			log.Println("all deployments complete")
			break
		}

		time.Sleep(time.Second * time.Duration(interval))
	}
}

func main() {
	cluster := flag.String("cluster", "", "ecs cluster where service lives")
	statusInterval := flag.Int("interval", 5, "Status check interval in seconds")
	profile := flag.String("profile", "", "AWS profile for deployment. If not stated, will use default")
	region := flag.String("region", "us-west-2", "AWS region that service lives")
	service := flag.String("service", "", "service to trigger deployment")
	flag.Parse()

	if *cluster == "" {
		log.Fatal("cluster is missing and required")
	}

	if *service == "" {
		log.Fatal("service is missing and required")
	}

	// Setup the session for ECS
	sessOptions := session.Options{
		Config: aws.Config{Region: aws.String(*region)},
	}
	if *profile != "" {
		sessOptions.Profile = *profile
	}

	sess := session.Must(session.NewSessionWithOptions(sessOptions))
	svc := ecs.New(sess)

	// Start deployment of service
	err := forceNewDeployment(svc, *cluster, *service)
	if err != nil {
		log.Fatal(err)
	}

	// track deployment status
	getDeploymentStatusTilComplete(svc, *cluster, *service, *statusInterval)
}
