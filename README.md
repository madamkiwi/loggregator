# Loggregator

[![CI Badge][ci-badge]][ci-pipeline]

### Logging in the Clouds

Loggregator is the user application logging subsystem of Cloud Foundry.

If you would to see our plans for future development, check out the [Loggregator Roadmap](https://docs.google.com/spreadsheets/d/1bM1bInPQeC2xLayLsFb0aBuD3_HFNfJj9mEJZygnuWo/edit#gid=0).
You can track what we're working on today in our [Pivotal Tracker Project](https://www.pivotaltracker.com/n/projects/993188).
If you need to contact us, you can [join our Slack channel](https://cloudfoundry.slack.com/messages/loggregator/).

### Table of Contents

* [Features](#features)
* [Usage](#usage)
* [Constraints](#constraints)
* [Architecture](#architecture)
* [Emitting Messages from other Cloud Foundry components](#emitting-messages-from-other-cloud-foundry-components)
* [Generating TLS Certificates](#generating-tls-certificates)
  * [Custom TLS Certificates](#custom-tls-certificate-generation)
* [Enabling TLS communication with etcd](#enabling-tls-between-loggregator-and-etcd)
* [Deploying via BOSH](#deploying-via-bosh)
* [Configuring the Firehose](#configuring-the-firehose)
* [Consuming Log and Metric Data](#consuming-log-and-metric-data)
* [Development](#development)
* [Loggregator as a separate release](#loggregator-as-a-separate-release)
* [Metrics generated by Loggregator](#metrics-generated-by-loggregator)
* [Logging message loss in Loggregator](#logging-message-loss-in-loggregator)
* [Loggregator does not guarantee logs delivered in order](#loggregator-does-not-guarantee-logs-delivered-in-order)
* [Workaround for Java multiline messages](#multi-line-java-message-workaround)
* [Log Message Size Constraints](#log-message-size-constraints)
* [Disable Metron UDP](#disable-metron-udp)
* [FAQ](https://github.com/cloudfoundry/loggregator/wiki/FAQ)
* Additional readmes
  * [Bosh HM metrics Forwarder](src/boshhmforwarder/README.md)
  * [DEA Logging Agent](src/deaagent/README.md)
  * [Doppler](src/doppler/README.md)
  * [LATS - Loggregator Acceptance Tests](src/lats/README.md)
  * [Metron](src/metron/README.md)
  * [Statsd-injector](src/statsd-injector/README.md)
    * [Statsd Go Client](src/statsd-injector/tools/statsdGoClient/README.md)
    * [Statsd Java Client](src/statsd-injector/tools/statsdJavaClient/README.md)
    * [Statsd Ruby Client](src/statsd-injector/tools/statsdRubyClient/README.md)
  * [Syslog Drain Binder](src/syslog_drain_binder/README.md)
  * [Doppler Benchmarking](src/tools/dopplerbenchmark/README.md)
  * [Metron Benchmarking](src/tools/metronbenchmark/README.md)
  * [Traffic Controller](src/trafficcontroller/README.md)
  * [Truncating Buffer](src/truncatingbuffer/README.md)

### Features

Loggregator allows users to:

1. Tail their application logs.
1. Dump a recent set of application logs (where recent is a configurable number of log packets).
1. Continually drain their application logs to 3rd party log archive and analysis services.
1. (Operators and administrators only) Access the firehose, which includes the combined stream of logs from all apps, plus metrics data from CF components.

### Usage

First, make sure you're using the new [golang based CF CLI](https://github.com/cloudfoundry/cli).  Once that's installed:

```
cf logs APP_NAME [--recent]
```
For example:

``` bash
$ cf logs private-app
Connected, tailing...
Oct 3 15:09:26 private-app App/0 STDERR This message is on stderr at 2013-10-03 22:09:26 +0000 for private-app instance 0
Oct 3 15:09:26 private-app App/0 STDERR 204.15.2.45, 10.10.2.148 - - [03/Oct/2013 22:09:26] "GET / HTTP/1.1" 200 81 0.0010
Oct 3 15:09:26 private-app App/0 This message is on stdout at 2013-10-03 22:09:26 +0000 for private-app instance 0
Oct 3 15:09:26 private-app App/0 STDERR This message is on stderr at 2013-10-03 22:09:26 +0000 for private-app instance 0
^C
```

### Constraints

1. Loggregator collects STDOUT & STDERR from applications.  This may require configuration on the developer's side.
  * **Warning**: The DEA logging agent expects an application to have open connections on **both** STDOUT and STDERR. Closing either of these (for example, by redirecting output to `/dev/null`) will be read by the logging agent as a misbehaving application, and it will disconnect from **all** sockets for that app.
1. A Loggregator outage must not affect the running application.
1. Loggregator gathers and stores logs in a best-effort manner.  While undesirable, losing the current buffer of application logs is acceptable.
1. The 3rd party drain API should mimic Heroku's in order to reduce integration effort for our partners.  The Heroku drain API is simply remote syslog over TCP.

### Architecture

Loggregator is composed of:

* **Sources**: Logging agents that run on the Cloud Foundry components.
* **Metron**: Metron agents are co-located with sources. They collect logs and forward them to:
* **Doppler**: Responsible for gathering logs from the **Metron agents**, storing them in temporary buffers, and forwarding logs to 3rd party syslog drains.
* **Traffic Controller**: Handles client requests for logs. Gathers and collates messages from all Doppler servers, and provides external API and message translation (as needed for legacy APIs).

Source agents emit the logging data as [protocol-buffers](https://github.com/google/protobuf), and the data stays in that format throughout the system.

![Loggregator Diagram](docs/loggregator.png)

In a redundant CloudFoundry setup, Loggregator can be configured to survive zone failures. Log messages from non-affected zones will still make it to the end user. On AWS, availability zones could be used as redundancy zones. The following is an example of a multi zone setup with two zones.

![Loggregator Diagram](docs/loggregator_multizone.png)

The role of Metron is to take traffic from the various emitter sources (dea, dea-logging-agent, router, etc) and route that traffic to one or more dopplers. In the current config we route this traffic to the dopplers in the same az. The traffic is randomly distributed across dopplers.

The role of Traffic Controller is to handle inbound HTTP and WebSocket requests for log data. It does this by proxying the request to all dopplers (regardless of AZ). Since an application can be deployed to multiple AZs, its logs can potentially end up on dopplers in multiple AZs. This is why the traffic controller will attempt to connect to dopplers in each AZ and will collate the data into a single stream for the web socket client.

The traffic controller itself is stateless; an incoming request can be handled by any instance in any AZ.

Traffic controllers also exposes a `firehose` web socket endpoint. Connecting to this endpoint establishes connections to all dopplers, and streams logs and metrics for all applications and CF components.

### Emitting Messages from other Cloud Foundry components

Cloud Foundry developers can easily add source clients to new CF components
that emit messages to the doppler.  Currently, there are libraries for
[Go](https://github.com/cloudfoundry/dropsonde/) and
[Ruby](https://github.com/cloudfoundry/loggregator_emitter). For usage
information, look at their respective READMEs.

### Generating TLS Certificates

To generate the Loggregator TLS certs and keys, run
`scripts/generate-loggregator-certs`.

#### Using bosh 2.0

If you are using the [new bosh cli](https://github.com/cloudfoundry/bosh-cli) you can generate certs using `--vars-store` flag.

```
bosh -e lite -d loggregator deploy templates/loggregator.yml
--vars-store=loggregator-vars.yml
```

#### Custom TLS Certificates

If you already have a CA, you'll need to create a certificate for each component.  Values for the components are:

##### Doppler
- common name: doppler
- extended key usage: serverAuth, clientAuth

##### TrafficController
- common name: trafficcontroller
- extended key usage: serverAuth, clientAuth

##### Metron
- common name: metron
- extended key usage: serverAuth, clientAuth

### Adding your TLS certificates

#### Doppler

| Property                       | Required | Description                                         |
|--------------------------------|----------|-----------------------------------------------------|
| `loggregator.tls.doppler.cert` | Yes      | Certificate used by doppler to communicate over TLS |
| `loggregator.tls.doppler.key`  | Yes      | Key used by doppler to communicate over TLS         |
| `loggregator.tls.ca_cert`      | Yes      | Certificate Authority used to sign the certificate  |

#### Traffic Controller

| Property                                 | Required | Description                                                    |
|------------------------------------------|----------|----------------------------------------------------------------|
| `loggregator.tls.trafficcontroller.cert` | Yes      | Certificate used by traffic controller to communicate over TLS |
| `loggregator.tls.trafficcontroller.key`  | Yes      | Key used by traffic controller to communicate over TLS         |
| `loggregator.tls.ca_cert`                | Yes      | Certificate Authority used to sign the certificate             |

#### Metron

| Property                      | Required | Description                                        |
|-------------------------------|----------|----------------------------------------------------|
| `loggregator.tls.metron.cert` | Yes      | Certificate used by metron to communicate over TLS |
| `loggregator.tls.metron.key`  | Yes      | Key used by metron to communicate over TLS         |
| `loggregator.tls.ca_cert`     | Yes      | Certificate Authority used to sign the certificate |

#### Example Manifest

```yaml
  loggregator:
    tls:
      ca_cert: |
        -----BEGIN CERTIFICATE-----
        LOGGREGATOR CA CERTIFICATE
        -----END CERTIFICATE-----
      doppler:
        cert: |
          -----BEGIN CERTIFICATE-----
          DOPPLER CERTIFICATE
          -----END CERTIFICATE-----
        key: |
          -----BEGIN RSA PRIVATE KEY-----
          DOPPLER KEY
          -----END RSA PRIVATE KEY-----
      trafficcontroller:
        cert: |
          -----BEGIN CERTIFICATE-----
          TRAFFIC CONTROLLER CERTIFICATE
          -----END CERTIFICATE-----
        key: |
          -----BEGIN RSA PRIVATE KEY-----
          TRAFFIC CONTROLLER KEY
          -----END RSA PRIVATE KEY-----
      metron:
        cert: |
          -----BEGIN CERTIFICATE-----
          METRON CERTIFICATE
          -----END CERTIFICATE-----
        key: |
          -----BEGIN RSA PRIVATE KEY-----
          METRON KEY
          -----END RSA PRIVATE KEY-----
```

### Enabling TLS between Loggregator and etcd

By default, metron, doppler, syslog_drain_binder, and loggregator_trafficcontroller all communicate with etcd over
http.  To enable TLS mutual auth to etcd, you'll need to generate certificates and update your manifest

**NOTE: Communicating with etcd over TLS is currently experimental, and will have an impact on other components within CloudFoundry. Enable it at your own discretion. The properties discussed below as well as their behavior might change in the future.**

#### etcd

Refer to [etcd-release's guide on Encrypting Traffic](https://github.com/cloudfoundry-incubator/etcd-release#encryption) for etcd's properties.

#### Metron Agent

| Property        | Required                              | Description                                     |
|-----------------|---------------------------------------|-------------------------------------------------|
| `loggregator.etcd.require_ssl` | No<br> Default: `false`                   | Enable ssl for all communcation with etcd |
| `loggregator.etcd.ca_cert`   | Yes if `loggregator.etcd.require_ssl` is set to `true` <br>Default: `""`              | PEM-encoded CA certificate            |
| `metron_agent.etcd.client_cert`   | Yes if `loggregator.etcd.require_ssl` is set to `true` <br>Default: `""`              | PEM-encoded client certificate            |
| `metron_agent.etcd.client_key`   | Yes if `loggregator.etcd.require_ssl` is set to `true` <br>Default: `""`              | PEM-encoded client key            |

#### Doppler

| Property        | Required                              | Description                                     |
|-----------------|---------------------------------------|-------------------------------------------------|
| `loggregator.etcd.require_ssl` | No<br> Default: `false`                   | Enable ssl for all communcation with etcd |
| `loggregator.etcd.ca_cert`   | Yes if `loggregator.etcd.require_ssl` is set to `true` <br>Default: `""`              | PEM-encoded CA certificate            |
| `doppler.etcd.client_cert`   | Yes if `loggregator.etcd.require_ssl` is set to `true` <br>Default: `""`              | PEM-encoded client certificate            |
| `doppler.etcd.client_key`   | Yes if `loggregator.etcd.require_ssl` is set to `true` <br>Default: `""`              | PEM-encoded client key            |

#### Traffic Controller

| Property        | Required                              | Description                                     |
|-----------------|---------------------------------------|-------------------------------------------------|
| `loggregator.etcd.require_ssl` | No<br> Default: `false`                   | Enable ssl for all communcation with etcd |
| `loggregator.etcd.ca_cert`   | Yes if `loggregator.etcd.require_ssl` is set to `true` <br>Default: `""`              | PEM-encoded CA certificate            |
| `traffic_controller.etcd.client_cert`   | Yes if `loggregator.etcd.require_ssl` is set to `true` <br>Default: `""`              | PEM-encoded client certificate            |
| `traffic_controller.etcd.client_key`   | Yes if `loggregator.etcd.require_ssl` is set to `true` <br>Default: `""`              | PEM-encoded client key            |

#### Syslog Drain Binder

| Property        | Required                              | Description                                     |
|-----------------|---------------------------------------|-------------------------------------------------|
| `loggregator.etcd.require_ssl` | No<br> Default: `false`                   | Enable ssl for all communcation with etcd |
| `loggregator.etcd.ca_cert`   | Yes if `loggregator.etcd.require_ssl` is set to `true` <br>Default: `""`              | PEM-encoded CA certificate            |
| `syslog_drain_binder.etcd.client_cert`   | Yes if `loggregator.etcd.require_ssl` is set to `true` <br>Default: `""`              | PEM-encoded client certificate            |
| `syslog_drain_binder.etcd.client_key`   | Yes if `loggregator.etcd.require_ssl` is set to `true` <br>Default: `""`              | PEM-encoded client key            |

An example manifest is given below:

```yaml
  loggregator:
    etcd:
      require_ssl: true
      ca_cert: |
        -----BEGIN CERTIFICATE-----
        ETCD CA CERTIFICATE
        -----END CERTIFICATE-----


  metron_agent:
    etcd:
      client_cert: |
        -----BEGIN CERTIFICATE-----
        METRON AGENT CERTIFICATE
        -----END CERTIFICATE-----
      client_key: |
        -----BEGIN RSA PRIVATE KEY-----
        METRON AGENT KEY
        -----END RSA PRIVATE KEY-----

  doppler:
    etcd:
      client_cert: |
        -----BEGIN CERTIFICATE-----
        DOPPLER CERTIFICATE
        -----END CERTIFICATE-----
      client_key: |
        -----BEGIN RSA PRIVATE KEY-----
        DOPPLER KEY
        -----END RSA PRIVATE KEY-----

  traffic_controller:
    etcd:
      client_cert: |
        -----BEGIN CERTIFICATE-----
        TRAFFIC CONTROLLER CERTIFICATE
        -----END CERTIFICATE-----
      client_key: |
        -----BEGIN RSA PRIVATE KEY-----
        TRAFFIC CONTROLLER KEY
        -----END RSA PRIVATE KEY-----

  syslog_drain_binder:
    etcd:
      client_cert: |
        -----BEGIN CERTIFICATE-----
        SYSLOG DRAIN BINDER CERTIFICATE
        -----END CERTIFICATE-----
      client_key: |
        -----BEGIN RSA PRIVATE KEY-----
        SYSLOG DRAIN BINDER KEY
        -----END RSA PRIVATE KEY-----
```

### Deploying via BOSH

Below are example snippets for deploying the DEA Logging Agent (source), Doppler, and Loggregator Traffic Controller via BOSH.

```yaml
jobs:
- name: dea_next
  templates:
  - name: dea_next
    release: cf
  - name: dea_logging_agent
    release: cf
  - name: metron_agent
    release: cf
  instances: 1
  resource_pool: dea
  networks:
  - name: cf1
    default:
    - dns
    - gateway
  properties:
    dea_next:
      zone: z1
    metron_agent:
      zone: z1
    networks:
      apps: cf1

- name: doppler_z1 # Add "doppler_zX" jobs if you have runners in zX
  templates:
  - name: doppler
    release: cf
  - name: syslog_drain_binder
    release: cf
  - name: metron_agent
    release: cf
  instances: 1  # Scale out as neccessary
  resource_pool: common
  networks:
  - name: cf1
  properties:
    doppler:
      zone: z1
    networks:
      apps: cf1

- name: loggregator_trafficcontroller_z1
  templates:
  - name: loggregator_trafficcontroller
    release: cf
  - name: metron_agent
    release: cf
  instances: 1  # Scale out as necessary
  resource_pool: common
  networks:
  - name: cf1
  properties:
    traffic_controller:
      zone: z1 # Denoting which one of the redundancy zones this traffic controller is servicing
    metron_agent:
      zone: z1
    networks:
      apps: cf1

properties:
  loggregator:
    servers:
      z1: # A list of loggregator servers for every redundancy zone
      - 10.10.16.14
    incoming_port: 3456
    outgoing_port: 8080

  loggregator_endpoint: # The end point sources will connect to
    shared_secret: loggregatorEndPointSharedSecret
    host: 10.10.16.16
    port: 3456
```

### Configuring the Firehose

The firehose feature includes the combined stream of logs from all apps, plus metrics data from CF components, and is intended to be used by operators and administrators.

Access to the firehose requires a user with the `doppler.firehose` scope.

The "cf" UAA client needs permission to grant this custom scope to users.
The configuration of the `uaa` job in Cloud Foundry [adds this scope by default](https://github.com/cloudfoundry/cf-release/blob/2a3d95417da3c59564daeecd754eb00862030cd6/jobs/uaa/templates/uaa.yml.erb#L111).
However, if your Cloud Foundry instance overrides the `properties.uaa.clients.cf` property in a stub, you need to add `doppler.firehose` to the scope list in the `properties.uaa.clients.cf.scope` property.

#### Configuring at deployment time (via deployment manifest)

In your deployment manifest, add

```yaml
properties:
  …
  uaa:
    …
    clients:
      …
      cf:
        scope: …,doppler.firehose
      …
      doppler:
        override: true
        authorities: uaa.resource
        secret: YOUR-DOPPLER-SECRET
```

(The `properties.uaa.clients.doppler.id` key should be populated automatically.) These are also set by default in [cf-properties.yml](https://github.com/cloudfoundry/cf-release/blob/master/templates/cf-properties.yml#L304-L307).

#### Adding scope to a running cluster (via `uaac`)

Before continuing, you should be familiar with the [`uaac` tool](http://docs.cloudfoundry.org/adminguide/uaa-user-management.html).

1. Ensure that doppler is a UAA client. If `uaac client get doppler` returns output like

  ```
    scope: uaa.none
    client_id: doppler
    resource_ids: none
    authorized_grant_types: authorization_code refresh_token
    authorities: uaa.resource
  ```

  then you're set.

  1. If it does not exist, run `uaac client add doppler --scope uaa.none --authorized_grant_types authorization_code,refresh_token --authorities uaa.resource` (and set its secret).
  1. If it exists but with incorrect properties, run `uaac client update doppler --scope uaa.none --authorized_grant_types "authorization_code refresh_token" --authorities uaa.resource`.
1. Grant firehose access to the `cf` client.
  1. Check the scopes assigned to `cf` with `uaac client get cf`, e.g.

    ```
      scope: cloud_controller.admin cloud_controller.read cloud_controller.write openid password.write scim.read scim.userids scim.write
      client_id: cf
      resource_ids: none
      authorized_grant_types: implicit password refresh_token
      access_token_validity: 600
      refresh_token_validity: 2592000
      authorities: uaa.none
      autoapprove: true
    ```

  1. Copy the existing scope and add `doppler.firehose`, then update the client

    ```
    uaac client update cf --scope "cloud_controller.admin cloud_controller.read cloud_controller.write openid password.write scim.read scim.userids scim.write doppler.firehose"
    ```

### Consuming log and metric data

The [NOAA Client](https://github.com/cloudfoundry/noaa) library, written in Golang, can be used by Go applications to consume app log data as well as the log + metrics firehose. If you wish to write your own client application using this library, please refer to the NOAA source and documentation.

Multiple subscribers may connect to the firehose endpoint, each with a unique subscription_id. Each subscriber (in practice, a pool of clients with a common subscription_id) receives the entire stream. For each subscription_id, all data will be distributed evenly among that subscriber's client pool.
In the event any subscription connections are lost, the messages will be rerouted to the remaining conenctions in the client pool.

### Development

The Cloud Foundry team uses GitHub and accepts contributions via [pull request](https://help.github.com/articles/using-pull-requests).

Follow these steps to make a contribution to any of our open source repositories:
1. Complete our CLA Agreement for [individuals](http://www.cloudfoundry.org/individualcontribution.pdf) or [corporations](http://www.cloudfoundry.org/corpcontribution.pdf)
1. Set your name and email

    ```
    git config --global user.name "Firstname Lastname"
    git config --global user.email "your_email@youremail.com"
    ```

1. Fork the repo (from `develop` branch to get the latest changes)
1. Make your changes on a topic branch, commit, and push to github and open a pull request against the `develop` branch.

Once your commits are approved by Travis CI and reviewed by the core team, they will be merged.

#### Go version support

Loggregator components are deployed to Cloud Foundry with Go 1.7.4.

#### OS X prerequisites

Use brew and do

    brew install go --cross-compile-all
    brew install direnv

Make sure you add the proper entry to load direnv into your shell. See `brew info direnv`
for details. To be safe, close the terminal window that you are using to make sure the
changes to your shell are applied.

#### Checkout

```
git clone https://github.com/cloudfoundry/loggregator --recursive
cd loggregator
direnv allow
```

Please run `scripts/install-git-hooks`

#### Additional go tools

Install go vet:

    go get golang.org/x/tools/cmd/vet

Install gosub
```
go get github.com/vito/gosub
```

#### Running tests

```
scripts/test
```

#### Running specific tests

```
go get github.com/onsi/ginkgo/ginkgo
cd src/metron # or any other component
ginkgo -r
```

#### Debugging

Doppler will dump information about the running goroutines to stdout if sent a `USR1` signal.

```
goroutine 1 [running]:
runtime/pprof.writeGoroutineStacks(0xc2000bc3f0, 0xc200000008, 0xc200000001, 0xca0000c2001fcfc0)
	/home/travis/.gvm/gos/go1.1.1/src/pkg/runtime/pprof/pprof.go:511 +0x7a
runtime/pprof.writeGoroutine(0xc2000bc3f0, 0xc200000008, 0x2, 0xca74765c960d5c8f, 0x40bbf7, ...)
	/home/travis/.gvm/gos/go1.1.1/src/pkg/runtime/pprof/pprof.go:500 +0x3a
....
```

### Loggregator as a separate release

There are cases when releases outside of Cloud Foundry would like to emit logs and metrics to the Loggregator system. In such cases we have instructions on using Loggregator as a separate release [here](docs/using_separate_release.md).

### Metrics generated by Loggregator

See this [list](docs/metric_descriptions.md).

### Logging message loss in Loggregator

One of the design priciples of Loggregator is to not put "back pressure" on apps. Therefore apps can log faster than the logging infrastructre can transport. In order to implement non-blocking, UDP was chosen as the transport protocol. Because of this, log message delivery is **not** guaranteed. If you wish to monitor log loss, there are metrics that can inform you.

There are two places loss can occur. The first place is inside the Diego cell where an app is logging so fast the local Metron can’t get the log messages off of the VM fast enough. To monitor that, you have to be on Diego, then the following works.

Metric #1 is `rep.logSenderTotalMessagesRead`, the number of messages sent by Diego apps. Metric #2 is `MetronAgent.DopplerForwarder.sentMessages`, the number of messages passed on by the Metron. #1-#2 will tell you the loss within the cell. You can monitor it on a cell by cell basis, or across all VMs based on the tags.

The other place for loss is between the Metron and Dopplers. To monitor that (and this works on DEAs too), metric #1 is `MetronAgent.DopplerForwarder.sentMessages` (same as #1 above) and metric #2 is `DopplerServer.udpListener.receivedMessageCount`. #1-#2 gives you loss between Metron and Doppler.

Using a graphing utility to display these absolute counter values as "rate per second" can give a good visualization of any loss.

The Loggregator team as a plan to improve this. You can read about it in our [Reliable Logging Feature Narrative](https://docs.google.com/document/d/1AhNDiuZef0rlNB9aCbBO-0DbBz_bC6iixHaWdjtBIHY/edit?usp=sharing).

### Loggregator does not guarantee logs delivered in order

The Loggregator system is designed to scale horizontally in order to handle the different load profiles of logs and metrics that are generated by the apps and the system components. As the size of the deployment grows, this logging load gets larger than can be handled by a single VM chain. This scaling is done by adding Dopplers. The logging load coming out of different VMs varies greatly. These VMs range from DEA and Diego app runners to system components such as the Router and the Cloud Controller. In order to evenly distribute the load over the scaled Dopplers, the Metrons distribute their log streams over the different Dopplers. The algorithm used is that for each log line (message), the Metron chooses a random Doppler to send it to. This means that each sequential log message from an application is routed through a different Doppler via a different network path. Because of the different paths, log lines take varying times to arrive at the destination nozzle or syslog and therefore the lines are not guaranteed to be in order.

### Multi-line Java message workaround

Having log lines out of order can be painful when monitoring Java stack traces. There is a work around involving modifing the Java log output to have your app reformat your stacktrace messages so any newline characters are replaced by a token; and then have your log parsing code replace that token with newline characters again. The following example is for the Java Logback library and ELK, but you may be able to use the same strategy for other Java log libraries and log aggregators.

With the Java Logback library you do this by adding `%replace(%xException){'\n','\u2028'}%nopex` to your logging config , and then use the following logstash conf.
```
# Replace the unicode newline character \u2028 with \n, which Kibana will display as a new line.
    mutate {
      gsub => [ "[@message]", '\u2028', "
"]
# ^^^ Seems that passing a string with an actual newline in it is the only way to make gsub work
    }
```
to replace the token with a regular newline again so it displays "properly" in Kibana.

### Log Message Size Constraints

Application logs are emitted via UDP messages from Diego/DEA to Metron. The maximum log message size from a DEA is ~60KiB (assuming overhead for the transporting envelope). Diego breaks up log messages greater than ~60KiB into multiple envelopes to mitigate this constraint.

### Disable Metron UDP

Metron accepts envelopes in two different protocols: UDP and gRPC. UDP is insecure and supports the [v1 envelopes](https://github.com/cloudfoundry/dropsonde-protocol). It will eventually be deprecated. gRPC supports [v2](https://github.com/cloudfoundry/loggregator-api). If you are running Metron and would like to not have UDP for security purposes you can set the [`metron_agent.disable_udp`](https://github.com/cloudfoundry/loggregator/blob/69e7f387ddb0eb04a358b7078bf37e6adcf69439/jobs/metron_agent/spec#L49-L51) property in the bosh manifest.

[ci-badge]: https://loggregator.ci.cf-app.com/api/v1/pipelines/loggregator/jobs/run-tests/badge
[ci-pipeline]: https://loggregator.ci.cf-app.com/
