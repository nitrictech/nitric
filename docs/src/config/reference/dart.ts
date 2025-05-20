import { NavGroup } from '../types'
import { SiDart } from 'react-icons/si'

export const DartReference: NavGroup = {
  title: 'Dart',
  icon: SiDart,
  items: [
    {
      title: 'Overview',
      href: '/reference/dart',
      breadcrumbRoot: true,
    },
    {
      title: 'APIs',
      items: [
        {
          title: 'api()',
          href: '/reference/dart/api/api',
        },
        {
          title: 'api.all()',
          href: '/reference/dart/api/api-all',
        },
        {
          title: 'api.get()',
          href: '/reference/dart/api/api-get',
        },
        {
          title: 'api.post()',
          href: '/reference/dart/api/api-post',
        },
        {
          title: 'api.put()',
          href: '/reference/dart/api/api-put',
        },
        {
          title: 'api.delete()',
          href: '/reference/dart/api/api-delete',
        },
        {
          title: 'api.patch()',
          href: '/reference/dart/api/api-patch',
        },
        {
          title: 'api.route()',
          href: '/reference/dart/api/api-route',
        },
        {
          title: 'api.route.all()',
          href: '/reference/dart/api/api-route-all',
        },
        {
          title: 'api.route.get()',
          href: '/reference/dart/api/api-route-get',
        },
        {
          title: 'api.route.post()',
          href: '/reference/dart/api/api-route-post',
        },
        {
          title: 'api.route.put()',
          href: '/reference/dart/api/api-route-put',
        },
        {
          title: 'api.route.delete()',
          href: '/reference/dart/api/api-route-delete',
        },
        {
          title: 'api.route.patch()',
          href: '/reference/dart/api/api-route-patch',
        },
      ],
    },
    {
      title: 'Batch',
      items: [
        {
          title: 'job()',
          href: '/reference/dart/batch/job',
        },
        {
          title: 'job.handler()',
          href: '/reference/dart/batch/job-handler',
        },
        {
          title: 'job.submit()',
          href: '/reference/dart/batch/job-submit',
        },
      ],
    },
    {
      title: 'Key Value Stores',
      items: [
        {
          title: 'kv()',
          href: '/reference/dart/keyvalue/keyvalue',
        },
        {
          title: 'kv.get()',
          href: '/reference/dart/keyvalue/keyvalue-get',
        },
        {
          title: 'kv.set()',
          href: '/reference/dart/keyvalue/keyvalue-set',
        },
        {
          title: 'kv.delete()',
          href: '/reference/dart/keyvalue/keyvalue-delete',
        },
        {
          title: 'kv.keys()',
          href: '/reference/dart/keyvalue/keyvalue-keys',
        },
      ],
    },
    {
      title: 'Topics',
      items: [
        {
          title: 'topic()',
          href: '/reference/dart/topic/topic',
        },
        {
          title: 'topic.publish()',
          href: '/reference/dart/topic/topic-publish',
        },
        {
          title: 'topic.subscribe()',
          href: '/reference/dart/topic/topic-subscribe',
        },
      ],
    },
    {
      title: 'Queues',
      items: [
        {
          title: 'queue()',
          href: '/reference/dart/queues/queue',
        },
        {
          title: 'queue.enqueue()',
          href: '/reference/dart/queues/queue-enqueue',
        },
        {
          title: 'queue.dequeue()',
          href: '/reference/dart/queues/queue-dequeue',
        },
      ],
    },
    {
      title: 'Secrets',
      items: [
        {
          title: 'secret()',
          href: '/reference/dart/secrets/secret',
        },
        {
          title: 'secret.put()',
          href: '/reference/dart/secrets/secret-put',
        },
        {
          title: 'secret.version()',
          href: '/reference/dart/secrets/secret-version',
        },
        {
          title: 'secret.latest()',
          href: '/reference/dart/secrets/secret-latest',
        },
        {
          title: 'secret.version.access()',
          href: '/reference/dart/secrets/secret-version-access',
        },
      ],
    },
    {
      title: 'Storage',
      items: [
        {
          title: 'bucket()',
          href: '/reference/dart/storage/bucket',
        },
        {
          title: 'bucket.on()',
          href: '/reference/dart/storage/bucket-on',
        },
        {
          title: 'bucket.file()',
          href: '/reference/dart/storage/bucket-file',
        },
        {
          title: 'bucket.files()',
          href: '/reference/dart/storage/bucket-files',
        },
        {
          title: 'file.exists()',
          href: '/reference/dart/storage/bucket-file-exists',
        },
        {
          title: 'file.read()',
          href: '/reference/dart/storage/bucket-file-read',
        },
        {
          title: 'file.write()',
          href: '/reference/dart/storage/bucket-file-write',
        },
        {
          title: 'file.delete()',
          href: '/reference/dart/storage/bucket-file-delete',
        },
        {
          title: 'file.getDownloadUrl()',
          href: '/reference/dart/storage/bucket-file-downloadurl',
        },
        {
          title: 'file.getUploadUrl()',
          href: '/reference/dart/storage/bucket-file-uploadurl',
        },
      ],
    },
    {
      title: 'SQL',
      items: [
        {
          title: 'sql()',
          href: '/reference/dart/sql/sql',
        },
        {
          title: 'sql.connectionString()',
          href: '/reference/dart/sql/sql-connection-string',
        },
      ],
    },
    {
      title: 'Schedules',
      items: [
        {
          title: 'schedule()',
          href: '/reference/dart/schedule/schedule',
        },
        {
          title: 'schedule.every()',
          href: '/reference/dart/schedule/schedule-every',
        },
        {
          title: 'schedule.cron()',
          href: '/reference/dart/schedule/schedule-cron',
        },
      ],
    },
    {
      title: 'Websockets',
      items: [
        {
          title: 'websocket()',
          href: '/reference/dart/websocket/websocket',
        },
        {
          title: 'websocket.on()',
          href: '/reference/dart/websocket/websocket-on',
        },
        {
          title: 'websocket.send()',
          href: '/reference/dart/websocket/websocket-send',
        },
        {
          title: 'websocket.close()',
          href: '/reference/dart/websocket/websocket-close',
        },
      ],
    },
  ],
}
