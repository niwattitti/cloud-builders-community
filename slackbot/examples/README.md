# slackbot examples

Edit the `cloudbuild.yaml` files in this directory to include your Slack webhook URL:

```yaml
steps:
- name: 'gcr.io/$PROJECT_ID/slackbot'
  args: [ '--title', '$_TITLE',
          '--icon', '$_ICON',
          '--tag', '$_TAG',
          '--build', '$BUILD_ID',
          '--webhook', '<Add your webhook URL here>' ]
- name: 'gcr.io/cloud-builders/docker'
  args: [ 'build', '.', '-f', 'Dockerfile-success']
substitutions:
  _TITLE: 'Release trigger'
  _ICON: ':cloudbuild:'
  _SLACK_WEBHOOK_URL: https://hooks.slack.com/services/xxxxxxxx/xxxxxxxxx/xxxxxxxxxxxxxxxx
  _TAG: 'test'
```

Three examples are provided:

* Run `gcloud builds submit . --config=cloudbuild-success.yaml --substitutions=_SLACK_WEBHOOK_URL="<slack web hook url>"` to generate a notification for a successful build
* Run `gcloud builds submit . --config=cloudbuild-failure.yaml --substitutions=_SLACK_WEBHOOK_URL="<slack web hook url>"` to generate a notification for a failed build
* Run `gcloud builds submit . --config=cloudbuild-timeout.yaml --substitutions=_SLACK_WEBHOOK_URL="<slack web hook url>"` to generate a notification for a build which times out.
