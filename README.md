## README

### About

This is a CSE291 Virtualization course project about live migration of containers. We extended the `runc` library and added the `runc  live-migrate` command for automatic cross-host live migration of containers.

### Structure

- `ExampleContainer`: This is a sample CNN model training container that can be used to illustrate the live migration functionality.
  
- `runc`: This is forked from the release-1.1 version of the runc repo. We added the `migrate.go`, `receiver.go`, and `transfer.go` to implement live migration.
  

### Usage

#### Container Setup

To build a docker container from the `ExampleContainer` folder and use `runc` to execute it, we first need to setup the docker image from the Dockerfile.

```bash
docker build -t ....
```

Based on the docker image, we can build the root file system for the container.

```bash
docker export ... > rootfs.tar
```

Unzip the rootfs.tar to directory rootfs

```bash
tar -xvf ...
```

Now you can run the container with

```bash
sudo runc run <container_ID>
```

#### Runc Binary Setup

Go to the runc directory and run the following commands, and then the runc version that supports `live-migrate` will be installed on your machine.

```bash
make
sudo make install
```

#### NFS Setup

On the source machine, set up the NFS server:

Give permission to the client for accessing the shared directory

On the destination, mount the NFS directory by the following command:

#### Testing

On the destination host, run .... to prepare for receiving:

On the source machine, run:

You can see that the container starts off from the point it was stopped on the source host.