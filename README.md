## 概要

本リポジトリは、Go Gin（API モード） を用いて開発した、 Twitter クローンアプリのバックエンド API サーバーです。

フロントエンド（Next.js）と分離した API 専用構成で設計・実装しており、 認証・投稿など、SNS に必要な主要機能を網羅しています。

## URL
https://twitter-nextjs-frontend.vercel.app/

## 使用技術

- Go 1.25.5
- Gin 1.12.0
- PostgreSQL 14（Render）
- sqlc 1.24.0 （SQLから型安全なGoコードを生成）
- Redis（Render managed Redis）
- Docker
- Air（Goアプリの開発効率向上）
- Gmail SMTP(Gmail送信)

## 機能一覧

- サインアップ・ログイン
- ユーザーアクティベーション(サインアップ時にGmailメール送信)
- ツイート
- ツイート詳細
- ツイート一覧(※ページネーション機能付き)
- ユーザー詳細
- いいね
- リツイート
- ブックマーク
- ブックマーク一覧(※ページネーション機能付き)
- フォロー
- フォロー一覧(※ページネーション機能付き)
- フォロワー一覧(※ページネーション機能付き)
- メッセージ機能(DM)
- 退会機能

## 工夫した点

- フロントエンド（Next.js）と完全分離した API 専用構成で設計
- RESTに沿ったエンドポイント設計
- 型安全性を意識した設計（request/response定義 + sqlc）
- N+1問題を考慮したクエリ最適化（JOIN / サブクエリの使い分け）<br>Qiita記事: https://qiita.com/Rikuto-Iwashita/items/45107de8560c3238cf0e
- バリデーション・認証・セッション管理の実装（Redis使用）
- パスワードのハッシュ化によるセキュリティ強化

## 技術選定理由

- SNSアプリは多数のユーザーからの同時リクエストが発生するため、高い並行処理性能を持つGoを採用  

- Javaの学習経験があるため、同じ静的型付け言語であるGoを採用  
  → 既存の知識を活かし、スムーズに開発できると判断

- シンプルな構文で学習コストが低い点を考慮  
  → 個人開発でも開発スピードを維持しやすい

## ER図
<img width="1882" height="667" alt="image" src="https://github.com/user-attachments/assets/5dc11a3c-c8ad-4ff9-9c8c-747f1f0bede0" />

## 画面
#### サインアップ画面

<img width="1920" height="1080" alt="スクリーンショット 2026-05-05 21 16 31" src="https://github.com/user-attachments/assets/d66c068a-6ee5-4882-88da-daca754c2129" />

#### メール受信画面(Gmail)
<img width="715" height="334" alt="スクリーンショット 2026-05-05 21 34 14" src="https://github.com/user-attachments/assets/275ae217-f129-450d-90ee-12b833002910" />

#### メールのリンククリック後
<img width="1920" height="1080" alt="image" src="https://github.com/user-attachments/assets/0a67658d-e59f-43de-bbfd-46d6a778488f" />

#### ログイン画面
<img width="1920" height="1080" alt="スクリーンショット 2026-05-05 21 32 42" src="https://github.com/user-attachments/assets/25a53d1c-b6f1-40ac-bbbe-80c4fc7c91e3" />

#### ツイート画面
<img width="1920" height="1080" alt="スクリーンショット 2026-05-05 21 37 44" src="https://github.com/user-attachments/assets/0ebaabce-3612-49eb-8ed7-403c5703dfc4" />

#### ツイート覧
<img width="1920" height="1080" alt="スクリーンショット 2026-05-05 21 42 24" src="https://github.com/user-attachments/assets/aff7df7c-bd75-49fd-9080-988d5f52ad93" />

#### ツイート詳細
<img width="1920" height="1080" alt="スクリーンショット 2026-05-05 21 53 03" src="https://github.com/user-attachments/assets/28ef772c-9379-4d59-9b27-1b2e2e650f65" />

#### ユーザー詳細
<img width="1920" height="1080" alt="スクリーンショット 2026-05-05 21 49 36" src="https://github.com/user-attachments/assets/45a2daf5-810c-4900-9a73-d95be8b3e42a" />

#### いいね・リツイート・ブックマーク(動画)
https://github.com/user-attachments/assets/51090866-61c3-4132-8369-162a7b508353

#### ブックマーク一覧
<img width="1920" height="1080" alt="スクリーンショット 2026-05-05 21 52 04" src="https://github.com/user-attachments/assets/0386bdf1-7b8d-4a51-80da-e24e81ea198f" />

#### フォロー・フォロー解除(動画)
https://github.com/user-attachments/assets/c1697832-2723-404c-ade0-9acd3b587a66

#### フォロー一覧
<img width="1920" height="1080" alt="スクリーンショット 2026-05-05 22 42 57" src="https://github.com/user-attachments/assets/658815b6-321b-40bb-b96d-f068b2e87390" />

#### フォロワー一覧
<img width="1920" height="1080" alt="スクリーンショット 2026-05-05 22 42 57" src="https://github.com/user-attachments/assets/ba93ee4b-4584-4ed3-9f62-23551256d7ea" />

#### グループ作成
<img width="1920" height="1080" alt="スクリーンショット 2026-05-05 22 46 32" src="https://github.com/user-attachments/assets/058995d1-afb6-4239-8fa3-db30fcccf3bb" />

#### ユーザー追加
<img width="1920" height="1080" alt="スクリーンショット 2026-05-05 22 48 00" src="https://github.com/user-attachments/assets/4cb5b737-bc8c-4c14-89ca-3b68bfb1eadf" />

#### グループ一覧
<img width="1920" height="1080" alt="スクリーンショット 2026-05-05 22 49 17" src="https://github.com/user-attachments/assets/45bf4d6f-5833-4029-888c-a310d6289d08" />

#### メッセージ一覧
<img width="1920" height="1080" alt="スクリーンショット 2026-05-05 23 05 44" src="https://github.com/user-attachments/assets/ea73b83e-6c72-4c22-8f7d-a9b33f44c774" />

#### 退会(動画)
https://github.com/user-attachments/assets/e22fab7d-ab4c-4882-b5b7-10f9772d9ccd


