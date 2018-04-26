# `sd-card-backup`

A simple tool to back up SD cards. Just plug in SD cards for backup and run.

# Usage

Place a file at `~/.config/sd-card-backup/config.json` like this:

    {
      "destination_root": "/backup/path",
      "sd_card_mount_point": "/Volumes",
      "sd_card_names": [
        "KUBO",
        "NIXIE"
      ],
      "folder_mapping": [
        { "source": "DCIM",                      "destination": "DCIM"    },
        { "source": "PRIVATE/M4ROOT/CLIP",       "destination": "CLIP"    }
      ]
    }

`sd-card-backup` will iterate through any listed SD cards that are mounted and back up files sorted by `file type/month/SD card` as follows:


|Source|`/Volumes/KUBO/DCIM/103CANON/IMG_8868.CR2`|
|----|----|
|Destination|`/backup/path/Images/2018-04/KUBO/DCIM/103CANON/IMG_8868.CR2`|

|Source|`/Volumes/NIXIE/DCIM/101MSDCF/DSC07203.JPG`|
|----|----|
|Destination|`/backup/path/Images/2018-01/NIXIE/DCIM/101MSDCF/DSC07203.JPG`|

|Source|`/Volumes/NIXIE/PRIVATE/M4ROOT/CLIP/C0026.MP4`|
|----|----|
|Destination|`/backup/path/Videos/2018-02/NIXIE/CLIP/C0026.MP4`|
