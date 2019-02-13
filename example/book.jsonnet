local circuit_breaker = import 'circuit_breaker.libsonnet';
local routes = import 'routes.libsonnet';

{
  version: 1,
  dependencies: [
    {
      name: 'user',
      cluster_name: 'user-development',
      lb: 'user-app:8080',
      host_header: 'user-service',
      headers: {
        request_headers_to_add: [
          {
            header: {
              key: 'X-Test',
              value: 'abc',
            },
          },
        ],
      },
      tls: false,
      connect_timeout_ms: 250,
      circuit_breaker: circuit_breaker,
      routes: [
        {
          path: '/ping',
          timeout_ms: 100,
        },
        routes.root,
      ],
    },
    {
      name: 'ab-testing',
      cluster_name: 'ab-testing-development',
      sds: true,
      tls: false,
      connect_timeout_ms: 250,
      circuit_breaker: circuit_breaker,
      outlier_detection: {
        consecutive_5xx: 3,
      },
      health_checks: [
        {
          timeout: '3s',
          interval: '10s',
          unhealthy_threshold: 3,
          healthy_threshold: 3,
          no_traffic_interval: '30s',
          unhealthy_interval: '5s',
          unhealthy_edge_interval: '1s',
          healthy_edge_interval: '1s',
          event_log_path: '/dev/stderr',
          always_log_health_check_failures: true,
        },
      ],
      routes: [
        {
          path: '/grpc.health.v1.Health/Check',
          method: 'POST',
          timeout_ms: 3000,
          retry_policy: {
            retry_on: '5xx,connect-failure,refused-stream,cancelled,deadline-exceeded,resource-exhausted',
            num_retries: 3,
            per_try_timeout_ms: 700,
          },
        },
        routes.root,
      ],
    },
    {
      name: 'fault-user',
      cluster_name: 'fault-user-development',
      lb: 'user-app:8080',
      host_header: 'fault-user-service',
      tls: false,
      connect_timeout_ms: 250,
      circuit_breaker: circuit_breaker,
      routes: [routes.root],
      fault_filter_config: {
        abort: {
          http_status: 503,
          percentage: {
            numerator: 100,
            denominator: 'HUNDRED',
          },
        },
      },
    },
  ],
}
