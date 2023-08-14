#include <linux/init.h>
#include <linux/module.h>
#include <linux/cdev.h>
#include <linux/device.h>
#include <linux/kernel.h>
#include <linux/uaccess.h>
#include <linux/fs.h>
#include <linux/version.h>

#define MAX_DEV 5
#define BUFSIZE 100

static int max_dev=MAX_DEV;
module_param(max_dev,int,0660);

static int buffersize=BUFSIZE;
module_param(buffersize,int,0660);

static char *default_msg = "uninitialised";
module_param(default_msg,charp,0660);

char *buffer;

static int ptemplate_open(struct inode *inode, struct file *file);
static int ptemplate_release(struct inode *inode, struct file *file);
static long ptemplate_ioctl(struct file *file, unsigned int cmd, unsigned long arg);
static ssize_t ptemplate_read(struct file *file, char __user *buf, size_t count, loff_t *offset);
static ssize_t ptemplate_write(struct file *file, const char __user *buf, size_t count, loff_t *offset);

static const struct file_operations ptemplate_fops = {
    .owner      = THIS_MODULE,
    .open       = ptemplate_open,
    .release    = ptemplate_release,
    .unlocked_ioctl = ptemplate_ioctl,
    .read       = ptemplate_read,
    .write       = ptemplate_write
};

struct mychar_device_data {
    struct cdev cdev;
};

static int dev_major = 0;
static struct class *ptemplate_class = NULL;
static struct mychar_device_data ptemplate_data[MAX_DEV];

#if LINUX_VERSION_CODE >= KERNEL_VERSION(6,0,0)
static int ptemplate_uevent(const struct device *dev, struct kobj_uevent_env *env)
{
    add_uevent_var(env, "DEVMODE=%#o", 0666);
    return 0;
}
#else
static int ptemplate_uevent(struct device *dev, struct kobj_uevent_env *env)
{
    add_uevent_var(env, "DEVMODE=%#o", 0666);
    return 0;
}
#endif


static int __init ptemplate_init(void)
{
    int err, i;
    int msg_len;
    dev_t dev;

    // some defaults in case of insane parameters 
    if(max_dev < 1)
        max_dev = 1;

    if(buffersize < 1 || buffersize > PAGE_SIZE) 
        buffersize = BUFSIZE;

    //add an extra byte to the end of the buffer so it holds the amount of data required but 
    //theres still room to null terminate it
    buffersize += 1;

    pr_info("allocating %d bytes", max_dev * buffersize);

    //allocate a buffer of memory to back the devices
    //this is sized buffersize bytes for each device created as a single block
    //each device will then read from this block at an offset of buffersize*minor_number
    buffer = kzalloc(max_dev * buffersize, GFP_KERNEL);

    // copy the default string into the memory buffer at each devices offset
    msg_len = strlen(default_msg);
    if ( msg_len > buffersize){
	//*(default_msg+buffersize)=0;
	msg_len = buffersize-1;
    }

    for(i=0; i < max_dev;i++){
       strncpy((buffer+(i*buffersize)), default_msg , msg_len);
       *(buffer+(i*buffersize) + msg_len) = 0;
    } 

    err = alloc_chrdev_region(&dev, 0, max_dev, "ptemplate");

    dev_major = MAJOR(dev);

    // create a class in sysfs /sys/class/ptemplate
#if LINUX_VERSION_CODE >= KERNEL_VERSION(6,0,0)
    ptemplate_class = class_create("ptemplate");
#else
    ptemplate_class = class_create(THIS_MODULE, "ptemplate");
#endif
    ptemplate_class->dev_uevent = ptemplate_uevent;

    for (i = 0; i < max_dev; i++) {
        cdev_init(&ptemplate_data[i].cdev, &ptemplate_fops);
        ptemplate_data[i].cdev.owner = THIS_MODULE;

        cdev_add(&ptemplate_data[i].cdev, MKDEV(dev_major, i), 1);

        device_create(ptemplate_class, NULL, MKDEV(dev_major, i), NULL, "ptemplate-%d", i);
    }

    return 0;
}

static void __exit ptemplate_exit(void)
{
    int i;

    for (i = 0; i < max_dev; i++) {
        device_destroy(ptemplate_class, MKDEV(dev_major, i));
    }

    class_unregister(ptemplate_class);

    unregister_chrdev_region(MKDEV(dev_major, 0), MINORMASK);
    pr_info("ptemplate removed and cleanup");
}

static int ptemplate_open(struct inode *inode, struct file *file)
{
    printk("MYCHARDEV: Device open\n");
    return 0;
}

static int ptemplate_release(struct inode *inode, struct file *file)
{
    printk("MYCHARDEV: Device close\n");
    return 0;
}

static long ptemplate_ioctl(struct file *file, unsigned int cmd, unsigned long arg)
{
    printk("MYCHARDEV: Device ioctl\n");
    return 0;
}

static ssize_t ptemplate_read(struct file *file, char __user *ubuf, size_t count, loff_t *offset)
{
    int len=0;
    int minor = MINOR(file->f_path.dentry->d_inode->i_rdev);

    if(*offset > 0 || count < buffersize)
        return 0;

    len = strlen(buffer+minor*buffersize); 
    if(raw_copy_to_user(ubuf, buffer+minor*buffersize, len))
        return -EFAULT;

    *offset = len;
    return len;
}

static ssize_t ptemplate_write(struct file *file, const char __user *ubuf, size_t count, loff_t *offset)
{
    int minor = MINOR(file->f_path.dentry->d_inode->i_rdev);

    if(*offset > 0 || count >= buffersize)
        return -EFAULT;
    pr_info("writing to %d", minor);
    if(raw_copy_from_user(buffer+minor*buffersize, ubuf, count))
        return -EFAULT;

    //make sure the string we're storing is null terminated
    // or interesting string corruptions occur
    *(buffer+(minor*buffersize)+count) = 0;
    pr_info("buffer[%d]=%s", minor, buffer+minor*buffersize);
    pr_info("buffer %lu ==  %lu",count, (minor*buffersize)+count);
    *offset = count;
    return count;
}

MODULE_LICENSE("GPL");
MODULE_AUTHOR("chris procter");

module_init(ptemplate_init);
module_exit(ptemplate_exit);
