function isApiCommand(request) {
    return request.uri.startsWith('/_api')
}

function transformToQuery(request) {
    const host = request.headers['host'][0].value;
    const subDomain = host.replace("urls.10pines.dev", "").replace(".", "") || "NONE";
    let uri = request.uri;
    if (uri === "/") {
        uri = "/NONE"
    }
    request.uri = '/query/' + subDomain + uri;
    return request;
}

exports.handler = async (event) => {
    const request = event.Records[0].cf.request;
    if (isApiCommand(request)) {
        return request;
    }
    return transformToQuery(request);
};
