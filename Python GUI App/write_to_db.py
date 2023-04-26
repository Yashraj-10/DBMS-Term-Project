import time
log_file=open("metrics.log",'r')
i=1
while(1):
    time.sleep(3)
    metrics=log_file.readlines()
    for metric in metrics[i:i+56]:
        print(metric)
        # print(metric[0])
    i+=56