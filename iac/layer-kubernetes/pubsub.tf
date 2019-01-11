
// Raw topics

resource "google_pubsub_topic" "twitter-raw" {
  name = "twitter-raw"
}

resource "google_pubsub_subscription" "twitter-raw-sub" {
  name  = "twitter-raw-sub"
  topic = "${google_pubsub_topic.twitter-raw.name}"

  ack_deadline_seconds = 20
}

resource "google_pubsub_topic" "mastodon-raw" {
  name = "mastodon-raw"
}

resource "google_pubsub_subscription" "mastodon-raw-sub" {
  name  = "mastodon-raw-sub"
  topic = "${google_pubsub_topic.mastodon-raw.name}"

  ack_deadline_seconds = 20
}

