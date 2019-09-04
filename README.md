# StarlingX VirtualBox Installer

StarlingX VirtualBox Installer is a collection of tools designed for automating
installation of StarlingX images on VirtualBox instances according to official
instructions that can be found here (for an All-in-One system):

https://wiki.openstack.org/w/index.php?title=StarlingX/Containers/Installation&oldid=170746

## Description

This installer is actually a factory of installers.

Given the following inputs:
- the type of setup to be installed:
  + aiosx: all-in-one with one server (all services running on the same machine)
  + aiodx: all-in-one with two servers (every controller is also a compute)
  + standard: two controller nodes and two or more computes
  + storage: two controller nodes, two or more storage nodes, two or more computes
- configuration specific to the selected lab:
  + OAM network (floating IP, controller IP addresses, etc.)
  + VirtualBox NAT network
  + CPUs, memory and disk sizes for each node type
  + number of compute and storage nodes, etc.

The factory will output an archive containing shell scripts to:
- setup VirtualBox VMs and networking
- patch ISO with configuration needed by the automated installer
- automatically install a StarlingX release

This approach makes the install process more transparent because generated
shell scripts are readable and easy to update in case that's needed.

## Usage

To start the factory web server:

    stxlab
    
Will start an HTTP server on port 3000. The port number can be changed using
command line option `-port`. Main page contains 4 tabs, one for each setup
type: aiosx, aiodx, standard and storage. After updating configuration options
(already populated with defaults) user can click `Generate ... Installer`
button to download an archive containing installer scripts.

The web server is just a simple GUI front-end to access functionality provided
by installer-generator methods. The same functionality can be accessed via
CLI using binaries provided for each setup type: `aiosx`, `aiodx`, `standard`,
`storage`.

Each binary has a large set of command line options to fine tune VirtualBox
VMs and the installed StarlingX. All of them have default values so user can
run the commands without providing any option.

As with the web server, the result is an archive containing installer scripts.
This enables scenarios like:
- local deployment of install scripts from local factory:

      ./aiosx | tar -C /path/to/vbox/workspace -x

 - remote deployment of install scripts from local factory:

      ./aiosx | ssh user@vbox-machine tar -C /path/to/vbox/workspace -x
      
- local deployment of install scripts from remote factory:

      ssh user@installer-factory /path/to/aiosx | tar -C /path/to/vbox/workspace -x
      
The generated archive contains:
- `vbox-setup.sh`
- `prepare-bootimage.sh`
- `install.sh`
- `vmctl.sh`

Install steps are the following:

1. unpack installer archive. Assuming an all-in-one single node setup the
   instructions are:

       aiosx | tar -x
       cd aiosx
       
2. create VirtualBox VMs:

       ./vbox-setup.sh /path/to/starlingx/iso
       
   This will remove already existing VMs named `<setup-type>-controller-0`, etc.
   and create a new set of VMs according to specification then group them under
   `<setup-type>` (where setup-type can be aiosx, aiodx, standard, storage).
   
   Then vbox-setup.sh calls prepare-bootimage.sh to patch the provided ISO
   image with the lab configuration and to make the automated install process
   independent of console access via serial port. The patched ISO goes through
   CentOS Anaconda steps then reboots to a state where it is accessible via
   SSH.
   
3. (optional) provide OpenStack charts:

       cp /path/to/stx-openstack/chart.tgz ./stx-openstack.tgz
       
   Note that the name of the file must match exactly `stx-openstack.tgz`.
   Installer applies platform-integ-apps then exits if it can't find this file.
   
4. run automated installer:

       ./install.sh
      
   Which will go through the steps of running ansible, adding, configuring and
   unlocking additional nodes, waiting for platform-integ-apps then applying
   stx-openstack (when available).
   
   If the install procedure crashes at some point and you need to resume it
   from a specific step there's a shell here-doc marker `SKIP_SECTION` that
   can be moved to include the section of the file you need to skip (meaning
   all steps that were already executed).
   
`vmctl.sh` is provided in case you need to snapshot or restore the entire
setup. It pauses all the VMs before taking a snapshot and restores them to
a paused state then quickly resumes all of them to minimize clock skew between
VMs.

`vbox-setup.sh` creates an `env.sh` that can be sourced to enable shell
aliases for serial console access (they map to `socat` commands). For example:

    source ./env.sh
     aiosx-controller-0-serial

`install.sh` creates a local `ssh.config` file that can be used to
access installed nodes, for example `ssh -F ssh.config controller-0`
      
## Build

Clone this repository and run `make` to build it:

    make

Will create in `./bin` binaries that are self contained (include also template
files, thanks to `packr`) and can be moved/deployed somewhere else.
