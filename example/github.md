```json
{
  "name": "Create GitHub Issue",
  "steps": [
    {
      "type": "http_request",
      "method": "POST",
      "URL": "https://api.github.com/repos/{{owner}}/{{repo}}/issues",
      "headers": {
        "Authorization": "Bearer $GITHUB_PAT",
        "Accept": "application/vnd.github+json"
      },
      "body": {
        "title": "{{title}}",
        "body": "{{body}}"
      },
      "timeout": 5000,
      "retries": 3
    }
  ]
}
```



```json
{
    "body": "Created by Fline",
    "owner": "shubh-man007",
    "repo": "Fline",
    "result": {
        "active_lock_reason": null,
        "assignee": null,
        "assignees": [],
        "author_association": "OWNER",
        "body": "Created by Fline",
        "closed_at": null,
        "closed_by": null,
        "comments": 0,
        "comments_url": "https://api.github.com/repos/shubh-man007/Fline/issues/1/comments",
        "created_at": "2026-04-23T18:14:35Z",
        "events_url": "https://api.github.com/repos/shubh-man007/Fline/issues/1/events",
        "html_url": "https://github.com/shubh-man007/Fline/issues/1",
        "id": 4317932894,
        "issue_dependencies_summary": {
            "blocked_by": 0,
            "blocking": 0,
            "total_blocked_by": 0,
            "total_blocking": 0
        },
        "labels": [],
        "labels_url": "https://api.github.com/repos/shubh-man007/Fline/issues/1/labels{/name}",
        "locked": false,
        "milestone": null,
        "node_id": "I_kwDOSBrmPs8AAAABAV5tXg",
        "number": 1,
        "performed_via_github_app": null,
        "pinned_comment": null,
        "reactions": {
            "+1": 0,
            "-1": 0,
            "confused": 0,
            "eyes": 0,
            "heart": 0,
            "hooray": 0,
            "laugh": 0,
            "rocket": 0,
            "total_count": 0,
            "url": "https://api.github.com/repos/shubh-man007/Fline/issues/1/reactions"
        },
        "repository_url": "https://api.github.com/repos/shubh-man007/Fline",
        "state": "open",
        "state_reason": null,
        "sub_issues_summary": {
            "completed": 0,
            "percent_completed": 0,
            "total": 0
        },
        "timeline_url": "https://api.github.com/repos/shubh-man007/Fline/issues/1/timeline",
        "title": "Testing workflow",
        "updated_at": "2026-04-23T18:14:35Z",
        "url": "https://api.github.com/repos/shubh-man007/Fline/issues/1",
        "user": {
            "avatar_url": "https://avatars.githubusercontent.com/u/163862265?v=4",
            "events_url": "https://api.github.com/users/shubh-man007/events{/privacy}",
            "followers_url": "https://api.github.com/users/shubh-man007/followers",
            "following_url": "https://api.github.com/users/shubh-man007/following{/other_user}",
            "gists_url": "https://api.github.com/users/shubh-man007/gists{/gist_id}",
            "gravatar_id": "",
            "html_url": "https://github.com/shubh-man007",
            "id": 163862265,
            "login": "shubh-man007",
            "node_id": "U_kgDOCcRW-Q",
            "organizations_url": "https://api.github.com/users/shubh-man007/orgs",
            "received_events_url": "https://api.github.com/users/shubh-man007/received_events",
            "repos_url": "https://api.github.com/users/shubh-man007/repos",
            "site_admin": false,
            "starred_url": "https://api.github.com/users/shubh-man007/starred{/owner}{/repo}",
            "subscriptions_url": "https://api.github.com/users/shubh-man007/subscriptions",
            "type": "User",
            "url": "https://api.github.com/users/shubh-man007",
            "user_view_type": "public"
        }
    },
    "title": "Testing workflow"
}
```