"""
Generated by Erda Migrator.
Please implement the function entry, and add it to the list entries.
"""

import django.db.models
import json
import os
import requests


class PsGroupProjectsQuota(django.db.models.Model):
    """
    generated by erda-cli
    """

    id = django.db.models.BigIntegerField()
    created_at = django.db.models.DateTimeField(auto_now=True)
    updated_at = django.db.models.DateTimeField(auto_now=True, auto_now_add=True)
    project_id = django.db.models.BigIntegerField()
    project_name = django.db.models.CharField()
    prod_cluster_name = django.db.models.CharField()
    staging_cluster_name = django.db.models.CharField()
    test_cluster_name = django.db.models.CharField()
    dev_cluster_name = django.db.models.CharField()
    prod_cpu_quota = django.db.models.BigIntegerField()
    prod_mem_quota = django.db.models.BigIntegerField()
    staging_cpu_quota = django.db.models.BigIntegerField()
    staging_mem_quota = django.db.models.BigIntegerField()
    test_cpu_quota = django.db.models.BigIntegerField()
    test_mem_quota = django.db.models.BigIntegerField()
    dev_cpu_quota = django.db.models.BigIntegerField()
    dev_mem_quota = django.db.models.BigIntegerField()
    creator_id = django.db.models.BigIntegerField()
    updater_id = django.db.models.BigIntegerField()

    class Meta:
        db_table = "ps_group_projects_quota"


class PsGroupProjects(django.db.models.Model):
    """
    generated by erda-cli
    """

    id = django.db.models.BigIntegerField()
    name = django.db.models.CharField()
    display_name = django.db.models.CharField()
    logo = django.db.models.CharField()
    desc = django.db.models.CharField()
    cluster_config = django.db.models.CharField()
    cpu_quota = django.db.models.DecimalField()
    mem_quota = django.db.models.DecimalField()
    creator = django.db.models.CharField()
    org_id = django.db.models.BigIntegerField()
    version = django.db.models.CharField()
    created_at = django.db.models.DateTimeField()
    updated_at = django.db.models.DateTimeField()
    dd_hook = django.db.models.TextField()
    email = django.db.models.TextField()
    functions = django.db.models.CharField()
    active_time = django.db.models.DateTimeField()
    rollback_config = django.db.models.CharField()
    enable_ns = django.db.models.BooleanField()
    is_public = django.db.models.BooleanField()
    type = django.db.models.CharField()

    class Meta:
        db_table = "ps_group_projects"


class SPodInfo(django.db.models.Model):
    """
    generated by erda-cli
    """

    id = django.db.models.BigIntegerField()
    created_at = django.db.models.DateTimeField()
    updated_at = django.db.models.DateTimeField()
    cluster = django.db.models.CharField()
    namespace = django.db.models.CharField()
    name = django.db.models.CharField()
    org_name = django.db.models.CharField()
    org_id = django.db.models.CharField()
    project_name = django.db.models.CharField()
    project_id = django.db.models.CharField()
    application_name = django.db.models.CharField()
    application_id = django.db.models.CharField()
    runtime_name = django.db.models.CharField()
    runtime_id = django.db.models.CharField()
    service_name = django.db.models.CharField()
    workspace = django.db.models.CharField()
    service_type = django.db.models.CharField()
    addon_id = django.db.models.CharField()
    uid = django.db.models.CharField()
    k8s_namespace = django.db.models.CharField()
    pod_name = django.db.models.CharField()
    phase = django.db.models.CharField()
    message = django.db.models.CharField()
    pod_ip = django.db.models.CharField()
    host_ip = django.db.models.CharField()
    started_at = django.db.models.DateTimeField()
    cpu_request = django.db.models.FloatField()
    mem_request = django.db.models.BigIntegerField()
    cpu_limit = django.db.models.FloatField()
    mem_limit = django.db.models.BigIntegerField()

    class Meta:
        db_table = "s_pod_info"


class ProjectRequestRecord:
    """
    用以计算一个项目各环境下的 quota
    """

    def __init__(self, project_id: int, project_name: str):
        self.project_id = project_id
        self.project_name = project_name
        self.prod_cluster_name = ""
        self.staging_cluster_name = ""
        self.test_cluster_name = ""
        self.dev_cluster_name = ""

        self.prod_cpu_request = 0.0
        self.prod_mem_request = 0.0
        self.staging_cpu_request = 0.0
        self.staging_mem_request = 0.0
        self.test_cpu_request = 0.0
        self.test_mem_request = 0.0
        self.dev_cpu_request = 0.0
        self.dev_mem_request = 0.0


def get_es_url() -> str:
    """
    拼接 es url
    :return:
    """
    es_url = os.getenv("ES_URL")  # 普通安装升级时从 evnFrom ConfigMap 获取
    if es_url == "":
        es_url = os.getenv("PIPELINE_ES_URL")  # 集成时从流水线变量获取
    # 如果获取不到 ES_URL 立即退出
    if es_url == "":
        print("Quota Migration: can not find ES_URL from env 'ES_URL' or 'PIPELINE_ES_URL'")
        exit(1)
    return es_url + "/spot-machine_summary-full_cluster/_search"


def query_from_es(es_url):
    """
    从 es 查询每台机器的 request 情况
    :param es_url:
    :return:
    """
    payload = '{"query":{"bool":{"filter":{"range":{"timestamp":{"from":0,"include_lower":true,' \
              '"include_upper":true}}}}}} '
    headers = {
        "Content-Type": "application/json",
        "cache-control": "no-cache",
    }
    r = requests.post(url=es_url, headers=headers, data=payload)
    print("request clusters sources from ES, url: %s, response: %s".format(r.url, r.text))
    if r.status_code != 200:
        print("failed to request cluster resource from ES: status error. the response text:", r.text)
    print("the response of requesting ES: ", r.text)
    result = r.json()
    try:
        hits = result["hits"]
    except Exception as e:
        print(e)
        print("failed to request cluster resource from ES, response structure error: the key ['hits'] not found")
        exit(1)
    try:
        total = hits["total"]
    except Exception as e:
        print(e)
        print(
            "failed to request cluster resource from ES, response structure error: the key ['hits']['total'] not found")
        exit(1)
    try:
        hits = hits['hits']
    except Exception as e:
        print(e)
        print(
            "failed to request cluster resource form ES, response structure error, the key ['hits']['hits'] not found")
    return hits


def entry():
    """
    please implement this and add it to the list entries
    """
    # 如果没有任何项目, 则直接退出
    if PsGroupProjects.objects.all().count() == 0:
        print("没有任何项目, 跳过 quota 迁移")
        return

    # 查出所有项目的 cluster_config
    projects = dict()
    # 查询所有 s_pod_infos 记录（排除无效的记录）
    pod_infos = SPodInfo.objects.exclude(project_id__isnull=True) \
        .exclude(project_id="") \
        .exclude(workspace__isnull=True) \
        .exclude(workspace="") \
        .exclude(cpu_request__isnull=True) \
        .exclude(mem_request__isnull=True)

    # 记录每一个 cluster allocatable 总量
    cluster_cpu_allocatable_sum = dict()
    cluster_mem_allocatable_sum = dict()
    # 记录每一个 cluster request 总量
    cluster_cpu_request_sum = dict()
    cluster_mem_request_sum = dict()
    for pod_info in pod_infos:
        cluster_cpu_allocatable_sum[pod_info.cluster] = 0
        cluster_mem_allocatable_sum[pod_info.cluster] = 0
        cluster_cpu_request_sum[pod_info.cluster] = 0
        cluster_mem_request_sum[pod_info.cluster] = 0

    # 打印 s_pod_info 中查出的记录
    # print("s_pod_info 记录")
    # for item in pod_infos:
    #     print(json.dumps({"project_id": item.project_id, "project_name": item.project_name, "workspace": item.workspace,
    #                       "cpu_request": item.cpu_request, "mem_request": item.mem_request, "cluster": item.cluster},
    #                      indent=2))

    # 项目当前各环境的 request 情况乘以 1.3，作为各环境的 quota
    for pod_info in pod_infos:
        projects[pod_info.project_id] = ProjectRequestRecord(project_id=pod_info.project_id,
                                                             project_name=pod_info.project_name)
    for pod_info in pod_infos:
        if str(pod_info.workspace) == "prod":
            projects[pod_info.project_id].prod_cluster_name = pod_info.cluster
            projects[pod_info.project_id].prod_cpu_request += pod_info.cpu_request * 1.3 * 1000  # Core => MillCore
            projects[pod_info.project_id].prod_mem_request += pod_info.mem_request * 1.3 * 1024 * 1024  # Mib  => byte
        if str(pod_info.workspace) == "staging":
            projects[pod_info.project_id].staging_cluster_name = pod_info.cluster
            projects[pod_info.project_id].staging_cpu_request += pod_info.cpu_request * 1.3 * 1000
            projects[pod_info.project_id].staging_mem_request += pod_info.mem_request * 1.3 * 1024 * 1024
        if str(pod_info.workspace) == "test":
            projects[pod_info.project_id].test_cluster_name = pod_info.cluster
            projects[pod_info.project_id].test_cpu_request += pod_info.cpu_request * 1.3 * 1000
            projects[pod_info.project_id].test_mem_request += pod_info.mem_request * 1.3 * 1024 * 1024
        if str(pod_info.workspace) == "dev":
            projects[pod_info.project_id].dev_cluster_name = pod_info.cluster
            projects[pod_info.project_id].dev_cpu_request += pod_info.cpu_request * 1.3 * 1000
            projects[pod_info.project_id].dev_mem_request += pod_info.mem_request * 1.3 * 1024 * 1024
        # 累计各集群的 request 情况，后面按 request 比例分配 quota 时作为分母用到
        cluster_cpu_request_sum[pod_info.cluster] += pod_info.cpu_request * 1000
        cluster_mem_request_sum[pod_info.cluster] += pod_info.mem_request * 1024 * 1024

    # 打印目前的分配情况
    print("第一次分配情况")
    for key in projects:
        project = projects[key]
        d = {
            "project_id": project.project_id,
            "project_name": project.project_name,

            "prod_cluster_name": project.prod_cluster_name,
            "prod_cpu_request": project.prod_cpu_request,
            "prod_mem_request": project.prod_mem_request,

            "staging_cluster_name": project.staging_cluster_name,
            "staging_cpu_request": project.staging_cpu_request,
            "staging_mem_request": project.staging_mem_request,

            "test_cluster_name": project.test_cluster_name,
            "test_cpu_request": project.test_cpu_request,
            "test_mem_request": project.test_mem_request,

            "dev_cluster_name": project.dev_cluster_name,
            "dev_cpu_request": project.dev_cpu_request,
            "dev_mem_request": project.dev_mem_request,
        }
        print(d)

    # 从 ES 获取各个集群每台机器的 cpu mem allocatable 情况
    hits: list = query_from_es(get_es_url())
    for hit in hits:
        try:
            cluster_name = hit['_source']['tags']['cluster_name']
            cpu_request = hit['_source']['fields']['cpu_allocatable'] * 1000  # 单位 核*1000=毫核
            mem_request = hit['_source']['fields']['cpu_allocatable'] * 1024 * 1024  # 单位 Mib=>byte
            if cluster_name in cluster_cpu_allocatable_sum.keys():
                cluster_cpu_allocatable_sum[cluster_name] += cpu_request
            if cluster_name in cluster_mem_allocatable_sum.keys():
                cluster_mem_allocatable_sum[cluster_name] += mem_request
        except Exception as e:
            print("failed to get cluster name, cpu_request, mem_request", e)

    # 对每一个项目的四个环境的 quota 进行再分配: 从 allocatable 资源中取出一半, 按该环境的 request 对集群总 request 的比值分配
    for key in projects:
        # 生产
        cluster_name = projects[key].prod_cluster_name
        if cluster_name != "" and cluster_name in cluster_cpu_request_sum.keys():
            cpu_request_sum = cluster_cpu_request_sum[cluster_name]
            cpu_allocatable_sum = cluster_cpu_allocatable_sum[cluster_name] / 2
            mem_request_sum = cluster_mem_request_sum[cluster_name]
            mem_allocatable_sum = cluster_mem_allocatable_sum[cluster_name] / 2
            if cluster_name in cluster_cpu_allocatable_sum.keys() and \
                    cluster_name in cluster_cpu_request_sum.keys() and \
                    cpu_request_sum > 0:
                # 该环境 request 的值 / 该集群总的被 request 的值 * 该集群 allocatable 的值的一半
                projects[key].prod_cpu_request += projects[key].prod_cpu_request / cpu_request_sum * cpu_allocatable_sum
            if cluster_name in cluster_mem_allocatable_sum.keys() and \
                    cluster_name in cluster_mem_request_sum and \
                    mem_request_sum > 0:
                projects[key].prod_mem_request += projects[key].prod_mem_request / mem_request_sum * mem_allocatable_sum

        # 预发
        cluster_name = projects[key].staging_cluster_name
        if cluster_name != "" and cluster_name in cluster_cpu_request_sum.keys():
            cpu_request_sum = cluster_cpu_request_sum[cluster_name]
            cpu_allocatable_sum = cluster_cpu_allocatable_sum[cluster_name] / 2
            mem_request_sum = cluster_mem_request_sum[cluster_name]
            mem_allocatable_sum = cluster_mem_allocatable_sum[cluster_name] / 2
            if cluster_name in cluster_cpu_allocatable_sum.keys() and \
                    cluster_name in cluster_cpu_request_sum.keys() and \
                    cpu_request_sum > 0:
                projects[key].staging_cpu_request += projects[
                                                         key].staging_cpu_request / cpu_request_sum * cpu_allocatable_sum
            if cluster_name in cluster_mem_allocatable_sum.keys() and \
                    cluster_name in cluster_mem_request_sum and \
                    mem_request_sum > 0:
                projects[key].staging_mem_request += projects[
                                                         key].staging_mem_request / mem_request_sum * mem_allocatable_sum

        # 测试
        cluster_name = projects[key].test_cluster_name
        if cluster_name != "" and cluster_name in cluster_cpu_request_sum.keys():
            cpu_request_sum = cluster_cpu_request_sum[cluster_name]
            cpu_allocatable_sum = cluster_cpu_allocatable_sum[cluster_name] / 2
            mem_request_sum = cluster_mem_request_sum[cluster_name]
            mem_allocatable_sum = cluster_mem_allocatable_sum[cluster_name] / 2
            if cluster_name in cluster_cpu_allocatable_sum.keys() and \
                    cluster_name in cluster_cpu_request_sum.keys() and \
                    cpu_request_sum > 0:
                projects[key].test_cpu_request += projects[key].test_cpu_request / cpu_request_sum * cpu_allocatable_sum
            if cluster_name in cluster_mem_allocatable_sum.keys() and \
                    cluster_name in cluster_mem_request_sum and \
                    mem_request_sum > 0:
                projects[key].test_mem_request += projects[key].test_mem_request / mem_request_sum * mem_allocatable_sum

        # 开发
        cluster_name = projects[key].dev_cluster_name
        if cluster_name != "" and cluster_name in cluster_cpu_request_sum.keys():
            cpu_request_sum = cluster_cpu_request_sum[cluster_name]
            cpu_allocatable_sum = cluster_cpu_allocatable_sum[cluster_name] / 2
            mem_request_sum = cluster_mem_request_sum[cluster_name]
            mem_allocatable_sum = cluster_mem_allocatable_sum[cluster_name] / 2
            if cluster_name in cluster_cpu_allocatable_sum.keys() and \
                    cluster_name in cluster_cpu_request_sum.keys() and \
                    cpu_request_sum > 0:
                projects[key].dev_cpu_request += projects[key].dev_cpu_request / cpu_request_sum * cpu_allocatable_sum
            if cluster_name in cluster_mem_allocatable_sum.keys() and \
                    cluster_name in cluster_mem_request_sum and \
                    mem_request_sum > 0:
                projects[key].dev_mem_request += projects[key].dev_mem_request / mem_request_sum * mem_allocatable_sum

    print("最终分配情况")
    records = []
    for key in projects:
        project = projects[key]
        # 如果 quota 记录已存在，则忽略
        if PsGroupProjectsQuota.objects.filter(project_id=int(project.project_id)).count() > 0:
            continue
        quota_record = PsGroupProjectsQuota()
        quota_record.project_id = project.project_id
        quota_record.project_name = project.project_name
        quota_record.prod_cluster_name = project.prod_cluster_name
        quota_record.staging_cluster_name = project.staging_cluster_name
        quota_record.test_cluster_name = project.test_cluster_name
        quota_record.dev_cluster_name = project.dev_cluster_name
        quota_record.prod_cpu_quota = int(project.prod_cpu_request)
        quota_record.prod_mem_quota = int(project.prod_mem_request)
        quota_record.staging_cpu_quota = int(project.staging_cpu_request)
        quota_record.staging_mem_quota = int(project.staging_mem_request)
        quota_record.test_cpu_quota = int(project.test_cpu_request)
        quota_record.test_mem_quota = int(project.test_mem_request)
        quota_record.dev_cpu_quota = int(project.dev_cpu_request)
        quota_record.dev_mem_quota = int(project.dev_mem_request)
        quota_record.creator_id = 0
        quota_record.updater_id = 0
        if quota_record.prod_cluster_name == "" or \
                quota_record.staging_cluster_name == "" or \
                quota_record.test_cluster_name == "" or \
                quota_record.dev_cluster_name == "":
            try:
                project_record = PsGroupProjects.objects.get(pk=int(quota_record.project_id))
                try:
                    cluster_config = json.loads(project_record.cluster_config)
                    quota_record.prod_cluster_name = cluster_config["PROD"]
                    quota_record.staging_cluster_name = cluster_config["STAGING"]
                    quota_record.test_cluster_name = cluster_config["TEST"]
                    quota_record.dev_cluster_name = cluster_config["DEV"]
                except Exception as e:
                    print("获取 cluster_config 失败，cluster_config: ", project_record.cluster_config,
                          "project_id: ", project_record.project_id,
                          "project_name: ", project_record.project_name,
                          "Exception: ", e)
            except Exception as e:
                print("查询项目信息失败", "project_id:", int(quota_record.project_id), "Exception:", e)

        d = {
            "project_id": quota_record.project_id,
            "project_name": quota_record.project_name,

            "prod_cluster_name": quota_record.prod_cluster_name,
            "prod_cpu_quota": quota_record.prod_cpu_quota,
            "prod_mem_quota": quota_record.prod_mem_quota,

            "staging_cluster_name": quota_record.staging_cluster_name,
            "staging_cpu_quota": quota_record.staging_cpu_quota,
            "staging_mem_quota": quota_record.staging_mem_quota,

            "test_cluster_name": quota_record.test_cluster_name,
            "test_cpu_quota": quota_record.test_cpu_quota,
            "test_mem_quota": quota_record.test_mem_quota,

            "dev_cluster_name": quota_record.dev_cluster_name,
            "dev_cpu_quota": quota_record.dev_cpu_quota,
            "dev_mem_quota": quota_record.dev_mem_quota,
        }
        print("最终分配情况:", d)
        records.append(quota_record)

    PsGroupProjectsQuota.objects.bulk_create(records)

    pass


entries: [callable] = [
    entry,
]