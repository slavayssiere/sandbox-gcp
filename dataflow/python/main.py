
"""A streaming word-counting workflow.
"""

from __future__ import absolute_import

import argparse
import logging
import json
import re

from past.builtins import unicode

import apache_beam as beam
import apache_beam.transforms.window as window
from apache_beam.io import WriteToText
from apache_beam.metrics import Metrics
from apache_beam.metrics.metric import MetricsFilter
from apache_beam.options.pipeline_options import PipelineOptions
from apache_beam.options.pipeline_options import SetupOptions
from apache_beam.options.pipeline_options import StandardOptions


class WordExtractingDoFn(beam.DoFn):
  """Parse each line of input text into words."""

  def __init__(self):
    self.words_counter = Metrics.counter(self.__class__, 'words')
    self.word_lengths_counter = Metrics.counter(self.__class__, 'word_lengths')
    self.word_lengths_dist = Metrics.distribution(
        self.__class__, 'word_len_dist')
    self.empty_line_counter = Metrics.counter(self.__class__, 'empty_lines')

  def process(self, element):
    """Returns an iterator over the words of this element.
    The element is a line of text.  If the line is blank, note that, too.
    Args:
      element: the element being processed
    Returns:
      The processed element.
    """
    message = json.loads(element)
    text_line = message['data'].strip()
    if not text_line:
      self.empty_line_counter.inc(1)
    words = re.findall(r'[\w\']+', text_line, re.UNICODE)
    for w in words:
      self.words_counter.inc()
      self.word_lengths_counter.inc(len(w))
      self.word_lengths_dist.update(len(w))
    return words



def run(argv=None):
  """Build and run the pipeline."""
  parser = argparse.ArgumentParser()
  parser.add_argument('--output',
                      dest='output',
                      required=True,
                      help='Output file to write results to.')
  group = parser.add_mutually_exclusive_group(required=True)
  group.add_argument(
      '--input_topic',
      help=('Input PubSub topic of the form '
            '"projects/<PROJECT>/topics/<TOPIC>".'))
  group.add_argument(
      '--input_subscription',
      help=('Input PubSub subscription of the form '
            '"projects/<PROJECT>/subscriptions/<SUBSCRIPTION>."'))
  known_args, pipeline_args = parser.parse_known_args(argv)

  # We use the save_main_session option because one or more DoFn's in this
  # workflow rely on global context (e.g., a module imported at module level).
  pipeline_options = PipelineOptions(pipeline_args)
  pipeline_options.view_as(SetupOptions).save_main_session = True
  pipeline_options.view_as(StandardOptions).streaming = True
  p = beam.Pipeline(options=pipeline_options)

  # Read from PubSub into a PCollection.
  if known_args.input_subscription:
    messages = (p
                | beam.io.ReadFromPubSub(
                    subscription=known_args.input_subscription)
                .with_output_types(bytes))
  else:
    messages = (p
                | beam.io.ReadFromPubSub(topic=known_args.input_topic)
                .with_output_types(bytes))

  lines = messages | 'decode' >> beam.Map(lambda x: x.decode('utf-8'))

  # Count the occurrences of each word.
  def count_ones(word_ones):
    (word, ones) = word_ones
    return (word, sum(ones))

  counts = (lines
            | 'split' >> (beam.ParDo(WordExtractingDoFn())
                          .with_output_types(unicode))
            | 'pair_with_one' >> beam.Map(lambda x: (x, 1))
            | beam.WindowInto(window.FixedWindows(15, 0))
            | 'group' >> beam.GroupByKey()
            | 'count' >> beam.Map(count_ones))

  # Format the counts into a PCollection of strings.
  def format_result(word_count):
    (word, count) = word_count
    return '%s: %d' % (word, count)

  output = (counts
            | 'format' >> beam.Map(format_result)
            | 'encode' >> beam.Map(lambda x: x.encode('utf-8'))
            .with_output_types(bytes))

  # Write to GCS.
  # pylint: disable=expression-not-assigned
  output | 'write' >> WriteToText(known_args.output)

  result = p.run()
  result.wait_until_finish()


if __name__ == '__main__':
  logging.getLogger().setLevel(logging.INFO)
  run()