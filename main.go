package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		sgArgs := &ec2.SecurityGroupArgs{
			Ingress: ec2.SecurityGroupIngressArray{
				ec2.SecurityGroupIngressArgs{
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(8080),
					ToPort:     pulumi.Int(8080),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
				ec2.SecurityGroupIngressArgs{
					Protocol:   pulumi.String("tcp"),
					FromPort:   pulumi.Int(22),
					ToPort:     pulumi.Int(22),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
			Egress: ec2.SecurityGroupEgressArray{
				ec2.SecurityGroupEgressArgs{
					Protocol:   pulumi.String("-1"),
					FromPort:   pulumi.Int(0),
					ToPort:     pulumi.Int(0),
					CidrBlocks: pulumi.StringArray{pulumi.String("0.0.0.0/0")},
				},
			},
		}

		sg, err := ec2.NewSecurityGroup(ctx, "jenkins-sg", sgArgs)
		if err != nil {
			return err
		}

		kp, err := ec2.NewKeyPair(ctx, "local-ssh", &ec2.KeyPairArgs{
			PublicKey: pulumi.String("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQDgO2U+oOqxYutMoUCMmteFM7irj1B4XvPEbimf+6mU6dL+3h8pdfPM+aQbTnri/qdIUZiHTvYeo9Vix/JpCfUtArRKmYhpNWFpviUbYWrCeZL4QAWCvFfnG49h5fTzUx/22vp5iWW6r1OKYepGg76bnadqor/SyqskEEM1Gziy7YPdVW0U1QnIeRS1y0Eal5HoCRgel+RJwI4Tee/s6Fv/rykKtsva7ICOJUG45NIiTF8W2w+fulvmz//eZaT88zJSTeQYNXXYihy3LcNjMutgE2F0+TUDUVHHa4U16wXxXm0iyuLllGOxaDqDfdTtHgNXcu4YPbnSYhqbr0lHRrygJ8vzo0SStyJs3fl/A7PvykfX0PAmyzwM9M733r0qLHToDdbkXZfu6yMzO34G765t+oO8jLhGZ/5eRKvhtqDRdnizybDxj+joAO3QY/lxaOFSp0MyVs/aHOOjAk25GHHhYK7jLmkEQDubvDEarGRi+RYyBCmWdOGuEum//emHzn0= fajar@Fajars-MacBook-Pro.local"),
		})
		if err != nil {
			return err
		}

		jenkinsServer, err := ec2.NewInstance(ctx, "jenkins-server", &ec2.InstanceArgs{
			InstanceType:        pulumi.String("t2.micro"),
			VpcSecurityGroupIds: pulumi.StringArray{sg.ID()},
			Ami:                 pulumi.String("ami-0daef888755f9c098"),
			KeyName:             kp.KeyName,
		})
		if err != nil {
			return err
		}

		fmt.Println("Public IP: ", jenkinsServer.PublicIp)
		fmt.Println("Public DNS: ", jenkinsServer.PublicDns)

		ctx.Export("publicIP", jenkinsServer.PublicIp)
		ctx.Export("publicDNS", jenkinsServer.PublicDns)

		return nil
	})
}
