# Trace
A golang based path tracer with a simple web interface.
Path tracing based on Peter Shirley's [Ray Tracing in One Weekend](https://raytracing.github.io/books/RayTracingInOneWeekend.html).
![Ray trace showcase](https://github.com/g-jessmuir/trace/blob/master/samples/gotrace.gif)

## About
After writing too much C in too little time I wanted to learn something new. I'd previously worked through the book in C++, and decided to have another go in go. To build on the original project, I paralellized the ray tracing process and added a simple web UI (pictured above) to show progress. [Gorilla](https://github.com/gorilla/websocket) was used to handle websockets.

I didn't deviate much from the program structure described in the book so I could move through it quickly and focus on enhancements, leaving behind some not-so-idiomatic pieces here and there.

## Enhancements
* Parallel processing for speedier rendering
* A simple web interface for configuring the number of goroutines to use, the number of samples to take per pixel, and the seed for the random number generator
* Periodic updates to show incremental progress in the interface

## Lessons learned
* Using channels to handle incremental updates was a joy coming from C
* I missed operator overloading the most throughout the project, as it allows programming math on custom types to be much more descriptive and succinct
* Front-end development is not and likely never will be for me
