# Deployment

Provides different ways to deploy RCC. Intended as references to copy/paste to
your own deployment setup as a foundation, not to run directly. Unless you're
me, in which case I just run them directly.

## Lambda

The [lambda deploy](./terraform_lambda) is the actual deployment I use, purely
for cost as lambdas are very cheap for low traffic volume.

The link to `https://rcc.evertras.com` is done in [another repo](https://github.com/Evertras/site/blob/main/terraform/subdomain_rcc.tf).

## ECS

The [ECS deploy](./terraform_ecs) uses ECS in Fargate to run multiple container
instances connected to an ELB. The output var gives a link to the ELB, which
could then be used to link it further if desired.

This is a little nicer for versioning sanity, as opposed to needing to upload
code to the lambda each time, but generally the lambda seems simpler for this
particular use case.
