# Fibertel Station Exporter
Prometheus Exporter for the Fibertel Gateway [Technicolor CGA4233TCH3](https://www.technicolor.com/node/8053)

## Usage
```
Usage of ./fibertel-station-exporter:
  -log.level string
    	Logging level (default "info")
  -show-metrics
    	Show available metrics and exit
  -version
    	Print version and exit
  -fibertel.station-username string
    	Username for login into the Fibertel gateway
  -fibertel.station-password string
    	Password for login into the Fibertel gateway
  -fibertel.station-url string
    	Fibertel station URL. For bridge mode this is 192.168.100.1 (note: Configure a route if using bridge mode) (default "https://192.168.0.1")
  -web.listen-address string
    	Address to listen on (default "[::]:9420")
  -web.telemetry-path string
    	Path under which to expose metrics (default "/metrics")
```

## Exported metrics
* `fibertel_station_login_success_bool`: 1 if the login was successful
* `fibertel_station_login_message_info`: Login message returned by the web interface
  - Labels: `message`
* `fibertel_station_user_info`: User name as returned by the web interface
  - Labels: `username`
* `fibertel_station_uid_info`: User id as returned by the web interface
  - Labels: `uid`
* `fibertel_station_default_password_bool`: 1 if the default password is in use
* `fibertel_station_downstream_central_frequency_hertz`: Central frequency in hertz
  - Labels: `id`, `channel_id`, `fft`, `channel_type`
* `fibertel_station_downstream_power_dBmV`: Power in dBmV
  - Labels: `id`, `channel_id`, `fft`, `channel_type`
* `fibertel_station_downstream_snr_dB`: SNR in dB
  - Labels: `id`, `channel_id`, `fft`, `channel_type`
* `fibertel_station_downstream_snr_dB`: SNR in dB
  - Labels: `id`, `channel_id`, `fft`, `channel_type`
* `fibertel_station_ofdm_downstream_start_frequency_hertz`: Start frequency
  - Labels: `id`, `channel_id_ofdm`, `fft`, `channel_type`
* `fibertel_station_ofdm_downstream_end_frequency_hertz`: End frequency
  - Labels: `id`, `channel_id_ofdm`, `fft`, `channel_type`
* `fibertel_station_ofdm_downstream_central_frequency_hertz`: Central frequency
  - Labels: `id`, `channel_id_ofdm`, `fft`, `channel_type`
* `fibertel_station_ofdm_downstream_bandwidth_hertz`: Bandwidth
  - Labels: `id`, `channel_id_ofdm`, `fft`, `channel_type`
* `fibertel_station_ofdm_downstream_power_dBmV`: Power
  - Labels: `id`, `channel_id_ofdm`, `fft`, `channel_type`
* `fibertel_station_ofdm_downstream_snr_dB`: SNR
  - Labels: `id`, `channel_id_ofdm`, `fft`, `channel_type`
* `fibertel_station_ofdm_downstream_locked_bool`: Locking status
  - Labels: `id`, `channel_id_ofdm`, `fft`, `channel_type`
* `fibertel_station_upstream_central_frequency_hertz`: Central frequency
  - Labels: `id`, `channel_id_up`, `fft`, `channel_type`
* `fibertel_station_upstream_power_dBmV`: Power
  - Labels: `id`, `channel_id_up`, `fft`, `channel_type`
* `fibertel_station_upstream_ranging_status_info`: Ranging status
  - Labels: `id`, `channel_id_up`, `fft`, `channel_type`, `status`
* `fibertel_station_logout_success_bool`: 1 if the logout was successful
* `fibertel_station_logout_message_info`: Logout message returned by the web interface
  - Labels: `message`