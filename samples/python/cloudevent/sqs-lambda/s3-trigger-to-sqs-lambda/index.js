const {SQSClient, SendMessageCommand} = require("@aws-sdk/client-sqs");

// TODO this can be in a different file #cleanup
function paramsTemplate() {
  return {
    DelaySeconds: 10,
    MessageBody: "new Invoice",
    QueueUrl: process.env.QUEUE_URL
  };
}

exports.handler = async (event, attr) => {
    const sqsClient = new SQSClient({ region: process.env.REGION });
    for (const { s3, awsRegion, eventSource } of event.Records) {
      try {
        if (s3.object.key.toLowerCase().includes("secret")) {
          var params = paramsTemplate();
          const body = JSON.stringify({
            bucket: s3.bucket.name,
            key: s3.object.key
          });
          params.MessageBody = body;
          const command = new SendMessageCommand(params);
          const data = await sqsClient.send(command);
          
          // TODO it is sending null and I do not know why
          return data.httpStatusCode || 201;
        } else {
          return 304;
        }
      } catch (error) {
        console.error(error);
      }

    }
    return `Successfully processed ${event.Records.length} objects.`;
};
