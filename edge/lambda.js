exports.handler = async (event) => {
    const request = event.Records[0].cf.request;
    const host = request.headers['host'][0].value;
    const subDomain = host.replace("urls.10pines.dev", "").replace(".", "") || "none"
    request.headers['SubDomain'] = [{
        key: 'SubDomain',
        value: subDomain
    }];
    if (request.uri === '/') {
        request.uri = '/_none'
    }
    return request;
};
