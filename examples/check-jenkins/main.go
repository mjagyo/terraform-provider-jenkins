package main

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/mjagyo/jenkins-client-go"
)

// orderResource is the resource implementation.
type orderResource struct {
	client *jenkins.Client
}

type jobModel struct {
	Description types.String `json:"description"`
	Name        types.String `json:"name"`
	URL         types.String `json:"url"`
	Buildable   types.Bool   `json:"buildable"`
	Color       types.String `json:"color"`
}

func main() {
	url := "http://localhost:8080" // Define a string variable
	username := "admin"
	token := "11b79b85aaaef0653b94b4903986906680"

	client, err := jenkins.NewClient(&url, &username, &token)

	if err != nil {
		return
	}

	r := &orderResource{client: client}

	r.UpdateJob()

	// fmt.Printf("%v ====== ", jobs)

	// if err != nil {
	// 	return
	// }

	// fmt.Printf("%v !!!! ", jobs)
}

func (d *orderResource) Read() {
	jobs, err := d.client.GetJobs(nil)

	// Map response body to model
	for _, job := range jobs.Jobs {
		coffeeState := jobModel{
			Name:        types.StringValue(job.Name),
			Description: types.StringValue(job.Description),
			URL:         types.StringValue(job.URL),
			Buildable:   types.BoolValue(job.Buildable),
			Color:       types.StringValue(job.Color),
		}

		fmt.Printf("111 %v --- ", coffeeState)
	}

	if err != nil {
		return
	}
}

func (d *orderResource) CreateSecret() {
	// for _, item := range plan.Items {
	// 	items = append(items, jenkins.OrderItem{
	// 		Coffee: jenkins.Coffee{
	// 			ID: int(item.Coffee.ID.ValueInt64()),
	// 		},
	// 		Quantity: int(item.Quantity.ValueInt64()),
	// 	})
	// }

	// Create new order
	payload := jenkins.CredentialRequest{
		Credentials: jenkins.Credential{
			Scope:       "GLOBAL",
			ID:          "dsadsax1111cczx",
			Username:    "maaanu",
			Password:    "baccr",
			Description: "lixxxnda",
			Class:       "com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl",
		},
	}

	err := d.client.CreateSecret(payload)

	if err != nil {
		return
	}
}

func (d *orderResource) UpdateSecret() {
	payload := jenkins.Credential{
		StaperClass: "com.cloudbees.plugins.credentials.impl.UsernamePasswordCredentialsImpl",
		ID:          "usernamefromtf",
		Description: "a test 111",
		Username:    "cxzcxz",
	}

	err := d.client.UpdateSecret(payload)

	if err != nil {
		return
	}
}

func (d *orderResource) DeleteSecret() {
	err := d.client.DeleteSecret("identification")

	if err != nil {
		return
	}
}

func (d *orderResource) CreateJob() {
	err := d.client.CreateJob("test", "job-config.xml")

	if err != nil {
		return
	}
}

func (d *orderResource) UpdateJob() {
	err := d.client.UpdateJob("fromterraform", "changed-job.xml")

	if err != nil {
		return
	}
}
