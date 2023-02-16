package main

import (
	"bytes"
	"context"
	"log"
	"os"

	"github.com/kdomanski/iso9660"
	"gopkg.in/yaml.v3"
)

func CloudInitMain(ctx context.Context) error {
	writer, err := iso9660.NewWriter()
	if err != nil {
		log.Fatalf("failed to create writer: %s", err)
	}
	defer writer.Cleanup()

	b, _ := yaml.Marshal(map[string]interface{}{
		"hostname":          "helloworld",
		"disable_root":      false,
		"preserve_hostname": false,
		"system_info": map[string]interface{}{
			"distro":      "ubuntu",
			"ssh_svcname": "ssh",
		},
		"default_user": []string{},
		"users": []map[string]interface{}{
			{"name": "root", "ssh_pwauth": true, "ssh_authorized_keys": []string{
				"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDPyp+d7AKdQaV8zoO+iPV8kTlXMCSpcD7BR0kPc/S8UAyj8p2uVQWPCpJNo6qkK0g+hGZ1W4J9MihBAOCvL9DAd+XbUJr4bkkc7es3AkQ8OgmP+PR2zGR76cDyZ64tdiyYanWICm3spHAf2yusaJl9ETEletwSSrHLu25aeXD06Hu1y3rxqF1jl95Mu9HuBG7mksIRrlYExtlytl48NeO+xNFlslcgmoomowbrppnjgAcO7MfgjPnOcYjnYOhea+7cYG/ZFCwwstJ9ZP3zX2HY+4D+tuQeTNZsdFpu/JhaYZEshLopF1hB183kzpp6zrv+NisyHOv51cjZa64LfSZyvW8+tBrH+2oUjF5pMB4BXkMBVjpZjhlpDPgkmAcNV9ViE3Rp+Ov19Q3BUI195kw9yjdwlBH8lmxVs+f4loXaOz1cwS87zcAUbTZP2FLzoxNim+V1mqYCtmIoghpx0aeDoqknYpIWxlMBzvGylwMWBDUDDplsV8XN5ZHMjnBdoOc= puppy@LENOVO-YOGA"},
			},
		},
		"growpart": map[string]interface{}{
			"mode":    "auto",
			"devices": []string{"/"},
		},
	})
	err = writer.AddFile(bytes.NewBuffer([]byte("")), "meta-data")
	if err != nil {
		log.Fatalf("meta-data: %s", err)
	}
	err = writer.AddFile(bytes.NewBuffer(append([]byte("#cloud-config\n"), b...)), "user-data")
	if err != nil {
		log.Fatalf("user-data: %s", err)
	}

	outputFile, err := os.OpenFile("output.iso", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("failed to create file: %s", err)
	}

	err = writer.WriteTo(outputFile, "cidata")
	if err != nil {
		log.Fatalf("failed to write ISO image: %s", err)
	}

	err = outputFile.Close()
	if err != nil {
		log.Fatalf("failed to close output file: %s", err)
	}
	return nil
}
