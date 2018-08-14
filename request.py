import argparse
import requests
from requests.auth import HTTPBasicAuth
from common.bounded_thread_pool.bounded_executor import BoundedExecutor

class Total:
    total = 0

def requestWeb(url):
    http_proxy  = "http://127.0.0.1:8000"
    https_proxy = "https://127.0.0.1:8000"
    ftp_proxy   = "ftp://127.0.0.1:8000"

    proxyDict = { 
                "http"  : http_proxy, 
                "https" : https_proxy, 
                "ftp"   : ftp_proxy
                }

    r = requests.get(url, proxies=proxyDict)
    # r = requests.get(url)
    with open('result_matrix', 'a') as result_file:
        result_file.write("Status code: {} \n".format(r.status_code))
        result_file.write("Length of response content: {} \n".format(len(r.content)))
        Total.total += 1
        result_file.write("Now we have: {} \n".format(Total.total))

def main():
    links = []
    links.append("https://amazon.com")
    links.append("https://google.com")
    links.append("http://21cn.com")
    links.append("http://4399.com")
    links.append("https://baidu.com")
    links.append("http://pornhub.com")
    links.append("https://youtube.com")
    links.append("https://yahoo.com")
    links.append("http://www.mia.com")
    links.append("http://www.prada.com/cn/zh/women/bags/jcr:content/par/product_grid.a.sortBy_0.html")
    executor = BoundedExecutor(20)
    for url in links:
        for _ in range(100):
            executor.submit(requestWeb, url)

# def load_input():
#     parser = argparse.ArgumentParser()

#     parser.add_argument(
#         '-i', '--ip', help='IP of the request website', required=True, type=str, dest='i')
#     return parser.parse_args()

if __name__ == '__main__':
    main()