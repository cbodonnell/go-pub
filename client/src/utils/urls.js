import qs from 'qs';
import { environment } from '../environment';

const timestamp = (url) => `${url}?timestamp=${Date.now()}`;
const andTimestamp = (url) => `${url}&timestamp=${Date.now()}`;

function validURL(str) {
    var pattern = new RegExp('^(https?:\\/\\/)?'+ // protocol
      '((([a-z\\d]([a-z\\d-]*[a-z\\d])*)\\.)+[a-z]{2,}|'+ // domain name
      '((\\d{1,3}\\.){3}\\d{1,3}))'+ // OR ip (v4) address
      '(\\:\\d+)?(\\/[-a-z\\d%_.~+]*)*'+ // port and path
      '(\\?[;&a-z\\d%_.~+=-]*)?'+ // query string
      '(\\#[-a-z\\d_]*)?$','i'); // fragment locator
    return !!pattern.test(str);
}

function getSearchParam(search, param) {
  const query = qs.parse(search, { ignoreQueryPrefix: true });
  return query[param];
}

function isLocal(url) {
  const baseURL = window.location.protocol + '//' + window.location.hostname;
  return (url.indexOf(baseURL) === 0);
}

function wrapIfRemote(url, type, typeArg='') {
  return isLocal(url) ? url : `${environment.REACT_APP_ACTIVITY_URL}/${type}?${typeArg ? `type=${typeArg}&` : ''}remote=${encodeURIComponent(url)}`;
}

function routeRemote(url, type, typeArg='') {
  return `/${type}?${typeArg ? `type=${typeArg}&` : ''}remote=${encodeURIComponent(url)}`;
}

function proxyUrl(url) {
  return `${environment.REACT_APP_PROXY_URL}/?url=${encodeURIComponent(url)}`;
}
  

export {
  timestamp,
  andTimestamp,
  validURL,
  getSearchParam,
  isLocal,
  wrapIfRemote,
  routeRemote,
  proxyUrl
}