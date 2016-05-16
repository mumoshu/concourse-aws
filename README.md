### Auto-scaling Concourse CI v1.0 on AWS with Terraform

## Usage

1. Create 1 VPC and 2 subnets in it

2. Set up required environment variables required by the wrapper script for terraform
   ```
   $ cat >> .envrc <<<'
   export AWS_ACCESS_KEY_ID=<YOUR ACCESS KEY>
   export AWS_SECRET_ACCESS_KEY=<YOUR SECRET ACCESS KEY>
   export CONCOURSE_IN_ACCESS_ALLOWED_CIDR="<YOUR_PUBLIC_IP>/32"
   export CONCOURSE_SUBNET_ID=<YOUR_SUBNET1_ID>
   export CONCOURSE_DB_SUBNET_IDS=<YOUR_SUBNET1_ID>,<YOUR_SUBNET2_ID>
   '
   ```

   Install [direnv](https://github.com/direnv/direnv) and allow it to read `.envrc` created in the previous step.

   ```
   $ direnv allow
   ```

3. The same for optional ones
   ```
   $ export CONCOURSE_WORKER_INSTANCE_PROFILE=<YOUR INSTANCE PROFILE NAME>
   ```

4. Run the following commands to build required AMIs and to provision a Concourse CI cluster
   ```
   $ ./build-ubuntu-ami.sh
   $ ./build-docker-ami.sh
   $ ./build-concourse-ami.sh
   $ ./terraform.sh get
   $ ./terraform.sh plan
   $ ./terraform.sh apply
   ```

5. Open your browser and confirm that the Concourse CI is running on AWS:
   ```
   # This will extract the public hostname for your load balancer from terraform output and open your default browser
   $ open http://$(terraform output | ruby -e 'puts STDIN.first.split(" = ").last')
   ```

6. Follow the Concourse CI tutorial and experiment as you like:
   ```
   $ export CONCOURSE_URL=http://$(terraform output | ruby -e 'puts STDIN.first.split(" = ").last')
   $ fly -t test login -c $CONCOURSE_URL
   $ fly -t test set-pipeline -p hello-world -c hello.yml
   $ fly -t test unpause-pipeline -p hello-world
   ```
   See http://concourse.ci/hello-world.html for more information and the `hello.yml` referenced in the above example.

7. Modify autoscaling groups' desired capacity to scale out/in webs or workers.

## Why did you actually created this?

[BOSH](https://github.com/cloudfoundry/bosh) looks [very promising to me according to what problems it solves](https://bosh.io/docs/problems.html).
However I was too lazy to learn it for now mainly because:

* I'm not going to use IaaS other than AWS for the time being
* learning it to JUST try Concourse CI might be too much in the short term though

## You may also find those projects useful

* [Concourse CI docker image](https://github.com/MeteoGroup/concourse-ci)
* [gregarcara/concourse-docker](https://github.com/gregarcara/concourse-docker)
* [jtarchie/concourse-docker-compose](https://github.com/jtarchie/concourse-docker-compose)
  * I wonder if I could run docker containers instead of concourse ci's standalone binaries using this
* Maybe more up-to-date than [starkandwayne/terraform-concourse](https://github.com/starkandwayne/terraform-concourse)
* [motevets/concourse-in-a-box](https://github.com/motevets/concourse-in-a-box) to quickly get concourse up-and-running on a single Ubuntu 14.04 EC2 instance