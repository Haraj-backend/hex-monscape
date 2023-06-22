
# Observability

Observability is a way to understand a system from the outside by asking questions about it without knowing its inner workings. It helps us troubleshoot and handle new problems and answer the question "Why is this happening?".

To ask these questions, the application must be instrumented with signals such as `traces`, `metrics`, and `logs`. OpenTelemetry is the tool used to instrument application code and make a system observable.

We only use `traces` for our usecase.

## Traces

Trace records the paths journey taken by requests through multi-service architectures and helps pinpoint performance problems in distributed systems. It improves visibility and makes debugging easier by breaking down what happens within a request.

A trace is made of one or more spans, with the first span representing the root span and each root span representing a request from start to finish. The spans underneath provide more context of what occurs during a request.

This is what a trace looks like:

![trace](./trace.png)

To understand how tracing in OpenTelemetry works, here are the main components:

1. Trace Provider

   A Tracer Provider is a factory for Tracers. It is initialized once in most applications and its lifecycle matches the application’s lifecycle. It also includes Resource and Exporter initialization. It is typically the first step in tracing with OpenTelemetry.

2. Tracer

    A Tracer creates spans containing more information about what is happening for a given operation, such as a request in a service. Tracers are created from Tracer Providers.

3. Tracer Exporter

    Trace Exporters send traces to a consumer. Consumer could be OpenTelemetry Collector, Jaeger, AWS X-Ray, etc.

4. Trace Context (Context Propagation)

    With Context Propagation, Spans can be correlated with each other and assembled into a trace, regardless of where Spans are generated.

    A Context is an object that contains information for sending and receiving services to correlate spans and associate them with the overall trace.

    Propagation moves Context between services and processes to assembled. It serializes or deserializes Span Context and provides relevant Trace information to be propagated from one service to another.

    The default format used in OpenTelemetry tracing is `W3C TraceContext`.

## Spans

A span represents a unit of work or operation. Spans are the building blocks of Traces.

In OpenTelemetry, it include the following information:

1. Name
2. Parent span ID (empty for root spans)
3. Start and End Timestamps
4. Span Context

   Span context is an immutable object on every span that contains the following:

   - The Trace ID representing the trace that the span is a part of
   - The span’s Span ID
   - Trace Flags, a binary encoding containing information about the trace
   - Trace State, a list of key-value pairs that can carry vendor-specific trace information

5. Attributes

   Attributes are key-value pairs that contain metadata that you can use to annotate a Span to carry information about the operation it is tracking.

6. Span Events

   A Span Event can be thought of as a structured log message (or annotation) on a Span, typically used to denote a meaningful, singular point in time during the Span’s duration.

7. Span Links

    Links exist to associate one span with one or more spans, implying a causal relationship.

    For example, in a distributed system where some operations are tracked by a trace, an additional operation may be queued for asynchronous execution. To associate the trace for the subsequent operations with the first trace, a span link is used. The last span from the first trace can be linked to the first span in the second trace. Links are optional but serve as a good way to associate trace spans with one another.

8. Span Status

   A status will be attached to a span. Typically, you will set a span status when there is a known error in the application code, such as an exception.

   Status will be tagged as one of the following values:

   - Unset
   - Ok
   - Error

9. Span Kind

   When a span is created, it is one of Client, Server, Internal, Producer, or Consumer. This span kind provides a hint to the tracing backend as to how the trace should be assembled. If not provided, the span kind is assumed to be internal.

   Span Kind can be one of the following:

   - `Client`: A client span represents a synchronous outgoing remote call such as an outgoing HTTP request or database call. Note that in this context, "synchronous" does not refer to async/await, but to the fact that it is not queued for later processing.

   - `Server`: A server span represents a synchronous incoming remote call such as an incoming HTTP request or remote procedure call.

   - `Internal`: Internal spans represent operations which do not cross a process boundary. Things like instrumenting a function call or an express middleware may use internal spans.

   - `Producer`: Producer spans represent the creation of a job which may be asynchronously processed later. It may be a remote job such as one inserted into a job queue or a local job handled by an event listener.

   - `Consumer`: Consumer spans represent the processing of a job created by a producer and may start long after the producer span has already ended.
