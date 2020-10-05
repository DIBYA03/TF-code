mount_efs:
  file.directory:
    - name: /mnt/efs
  cmd.run:
    - name: mount -t nfs -o nfsvers=4.1,rsize=1048576,wsize=1048576,hard,timeo=600,retrans=2,noresvport {{ salt['environ.get']('EFS_DNS_NAME') }}:/ /mnt/efs/
