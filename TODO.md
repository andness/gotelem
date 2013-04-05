## TODO ##

    - Update to Go 1.1, a few cases of passing around methods as funcs
      that we can clean up.

    - Add an all-time summarizer that doesn't keep a list of seen
      observations but instead just keeps running totals for
      AVG,COUNT,MIN,MAX,SUM.

    - Need to clean up the path handling in the http publisher

    - Documentation

    - Make it possible to add your own custom Receivers

    - Thread-safety

    - Having a fixed max observation count for the HTTPPublisher works
      fine for the sampled data series since they are evenly spaced, but
      the raw observations can be used for x,y plots and it is hard to set
      a fixed limit on that. For that use you'd rather want a timestamp
      filter.

    - Multiple summarizers with the same observer source can share
      underlying storage. If for example you have a 1 minute, 5 minute and
      30 minute window, the 5 and 1 minute ones can piggyback on the 30
      minute since they both use a subset of the 30 minute data.

    - Sampling can be driven by a single ticker so that we don't have to
      create one ticker gper observer. With a few hundred observers that's
      a lot of waste. Since the user creates the observer by passing in
      the sample interval, we can cache the tickers in a common global.
