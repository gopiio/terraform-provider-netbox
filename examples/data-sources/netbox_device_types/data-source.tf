# Get all device types from a specific manufacturer
data "netbox_device_types" "by_manufacturer" {
  filter {
    name  = "manufacturer"
    value = "Nokia"
  }
}

# Get device types matching a model name
data "netbox_device_types" "by_model" {
  filter {
    name  = "model"
    value = "7210 SAS-Sx 10/100GE"
  }
}

# Get device types using a regex on the model name
data "netbox_device_types" "by_regex" {
  name_regex = "^7210.*"
}

# Get device types with multiple filters
data "netbox_device_types" "filtered" {
  filter {
    name  = "manufacturer"
    value = "Nokia"
  }
  filter {
    name  = "u_height"
    value = "1"
  }
}

# Limit the number of results
data "netbox_device_types" "limited" {
  limit = 10
  filter {
    name  = "manufacturer"
    value = "Nokia"
  }
}
