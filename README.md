# http2lambda

Invoke AWS Lambda using a simple HTTP call.

A simple URL handler that takes AWS credentials, function name and parameters from URL
parameters. Makes it possible to create and distribute a URL that invokes a Lambda function without
using AWS SDK.

There is apparently a bug in the AWS client library which causes Lambda to see the payload as a
string instead of an object (this is to be expected as the library is very immature). The Lambda
function that receives these calls needs to do something like this:

```javascript
exports.handler = function(event, context) {
  console.log('Received event: ' + typeof(event) + ': ' + event);

  if ('string' === typeof(event)) {
    var buf = new Buffer(event, 'base64');
    event = JSON.parse(buf);
  }
  ...
}
```
