package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAWSELBAttachment(t *testing.T) {
	var conf elb.LoadBalancerDescription

	testCheckInstanceAttached := func(count int) resource.TestCheckFunc {
		return func(*terraform.State) error {
			if len(conf.Instances) != count {
				return fmt.Errorf("instance count does not match")
			}
			return nil
		}
	}

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "aws_elb.bar",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckAWSELBDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAWSELBAttachmentConfig1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSELBExists("aws_elb.bar", &conf),
					testCheckInstanceAttached(1),
				),
			},

			resource.TestStep{
				Config: testAccAWSELBAttachmentConfig2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSELBExists("aws_elb.bar", &conf),
					testCheckInstanceAttached(2),
				),
			},

			resource.TestStep{
				Config: testAccAWSELBAttachmentConfig3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSELBExists("aws_elb.bar", &conf),
					testCheckInstanceAttached(0),
				),
			},
		},
	})

}

const testAccAWSELBAttachmentConfig1 = `
resource "aws_elb" "bar" {
  availability_zones = ["us-west-2a", "us-west-2b", "us-west-2c"]

  listener {
    instance_port = 8000
    instance_protocol = "http"
    lb_port = 80
    lb_protocol = "http"
  }
}

resource "aws_instance" "foo" {
	# us-west-2
	ami = "ami-043a5034"
	instance_type = "t1.micro"
}

resource "aws_elb_attachment" "app" {
  elb      = "${aws_elb.bar.id}"
  instances = ["${aws_instance.foo.id}"]
}
`

const testAccAWSELBAttachmentConfig2 = `
resource "aws_elb" "bar" {
  availability_zones = ["us-west-2a", "us-west-2b", "us-west-2c"]

  listener {
    instance_port = 8000
    instance_protocol = "http"
    lb_port = 80
    lb_protocol = "http"
  }
}

resource "aws_instance" "foo" {
	# us-west-2
	ami = "ami-043a5034"
	instance_type = "t1.micro"
}

resource "aws_instance" "foo2" {
	# us-west-2
	ami = "ami-043a5034"
	instance_type = "t1.micro"
}

resource "aws_elb_attachment" "app" {
  elb      = "${aws_elb.bar.id}"
  instances = ["${aws_instance.foo.id}", "${aws_instance.foo2.id}"]
}
`

const testAccAWSELBAttachmentConfig3 = `
resource "aws_elb" "bar" {
  availability_zones = ["us-west-2a", "us-west-2b", "us-west-2c"]

  listener {
    instance_port = 8000
    instance_protocol = "http"
    lb_port = 80
    lb_protocol = "http"
  }
}
`