exports.handler = async (event) => {
    const request = event.Records[0].cf.request;

    const host = request.headers['host'][0].value;
    const subDomain = host.replace("urls.10pines.dev", "").replace(".", "") || "none";

    const uri = request.uri;

    request.uri = '/query/' + subDomain + uri;
    return request;
};
