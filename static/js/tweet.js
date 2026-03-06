const urlParams = new URLSearchParams(window.location.search);
const LIMIT = 10;
let currentOffset = parseInt(urlParams.get('offset')) || 0;
let totalCount = 0;
let currentApiUrl = '';

// いいね機能
const setupLikeButton = (btn, tweetId, countLabel) => {
  btn.addEventListener('click', async () => {
    const isLiked = btn.dataset.isLiked === 'true';
    const method = isLiked ? 'DELETE' : 'POST';

    const response = await fetch(`/api/tweets/${tweetId}/like`, {
      method: method,
    });
    if (!response.ok) {
      console.error('サーバーエラーが発生しました', response.status);
      return;
    }

    const data = await response.json();

    btn.dataset.isLiked = data.is_liked;
    countLabel.textContent = data.like_count;
    btn.textContent = data.is_liked ? '❤️' : '♡';
  });
};

// リツイート機能
const setupRetweetButton = (btn, tweetId, countLabel) => {
  btn.addEventListener('click', async () => {
    const isRetweeted = btn.dataset.isRetweeted === 'true';
    const method = isRetweeted ? 'DELETE' : 'POST';

    const response = await fetch(`/api/tweets/${tweetId}/retweet`, {
      method: method,
    });
    if (!response.ok) {
      console.error('サーバーエラーが発生しました', response.status);
      return;
    }

    const data = await response.json();

    btn.dataset.isRetweeted = data.is_retweeted;
    countLabel.textContent = data.retweet_count;
    btn.textContent = data.is_retweeted ? '♻️' : '♻︎';
  });
};

// ブックマーク機能
const setupBookmarkButton = (btn, tweetId) => {
  btn.addEventListener('click', async () => {
    const isBookmarked = btn.dataset.isBookmarked === 'true';
    const method = isBookmarked ? 'DELETE' : 'POST';

    const response = await fetch(`/api/tweets/${tweetId}/bookmark`, {
      method: method,
    });
    if (!response.ok) {
      console.error('サーバーエラーが発生しました', response.status);
      return;
    }

    const data = await response.json();

    btn.dataset.isBookmarked = data.is_bookmarked;
    btn.textContent = data.is_bookmarked ? '📕' : '📖';
  });
};

// ツイートカード作成
const createTweetCard = (tweet) => {
  // デバッグ用
  console.log(tweet);
  const tweetCard = document.createElement('div');
  tweetCard.className =
    'p-4 border-b border-gray-100 hover:bg-gray-50/50 transition cursor-pointer';

  tweetCard.innerHTML = `
    <div class="flex gap-3">
      <div class="w-10 h-10 rounded-full bg-gray-200 flex-shrink-0"></div>
      <div class="flex-1">
        <div class="flex items-center gap-1">
          <span class="font-bold text-[15px] hover:underline">User ID: ${tweet.user_id}</span>
        </div>
        <p class="text-[15px] leading-5 mt-1 whitespace-pre-wrap">${tweet.content}</p>
        <button class="js-like-btn" data-is-liked="${tweet.is_liked}">${tweet.is_liked ? '❤️' : '♡'}</button> 
        <span class="js-like-count">${tweet.like_count}</span>
        <button class="js-retweet-btn" data-is-retweeted="${tweet.is_retweeted}">${tweet.is_retweeted ? '♻️' : '♻︎'}</button>
        <span class="js-retweet-count">${tweet.retweet_count}</span>
        <button class="js-bookmark-btn" data-is-bookmarked="${tweet.is_bookmarked}">${tweet.is_bookmarked ? '📕' : '📖'}</button>
      </div>
    </div>`;

  const likeBtn = tweetCard.querySelector('.js-like-btn');
  const likeLabel = tweetCard.querySelector('.js-like-count');
  const retweetBtn = tweetCard.querySelector('.js-retweet-btn');
  const retweetLabel = tweetCard.querySelector('.js-retweet-count');
  const bookmarkBtn = tweetCard.querySelector('.js-bookmark-btn');

  setupLikeButton(likeBtn, tweet.id, likeLabel);
  setupRetweetButton(retweetBtn, tweet.id, retweetLabel);
  setupBookmarkButton(bookmarkBtn, tweet.id);

  return tweetCard;
};

// ページネーションの初期設定
const setupPagination = () => {
  document.getElementById('prev-btn')?.addEventListener('click', () => {
    if (currentOffset >= LIMIT) {
      loadTweets(currentOffset - LIMIT);
    }
  });

  document.getElementById('next-btn')?.addEventListener('click', () => {
    if (currentOffset + LIMIT < totalCount) {
      loadTweets(currentOffset + LIMIT);
    }
  });
};

const updateUI = () => {
  const pageInfo = document.getElementById('page-info');
  const prevBtn = document.getElementById('prev-btn');
  const nextBtn = document.getElementById('next-btn');

  // ガード句を追加。tweet.jsは他のHTMLファイルでも読み込むため
  if (!pageInfo || !prevBtn || !nextBtn) return;

  const currentPage = Math.floor(currentOffset / LIMIT) + 1;
  const maxPage = Math.ceil(totalCount / LIMIT) || 1;

  pageInfo.textContent = `${currentPage} / ${maxPage} ページ (全 ${totalCount} 件)`;
  prevBtn.disabled = currentOffset === 0;
  nextBtn.disabled = currentOffset + LIMIT >= totalCount;
};

// ツイート一覧表示
const loadTweets = async (offset = 0) => {
  try {
    currentOffset = offset;

    // user-detail の時は URL を書き換えない
    if (!window.location.pathname.includes('user-detail')) {
      const params = new URLSearchParams(window.location.search);
      params.set('limit', LIMIT);
      params.set('offset', currentOffset);
      const newUrl = `${window.location.pathname}?${params.toString()}`;
      window.history.pushState({ offset: currentOffset }, '', newUrl);
    }

    const separator = currentApiUrl.includes('?') ? '&' : '?';
    const response = await fetch(
      `${currentApiUrl}${separator}limit=${LIMIT}&offset=${currentOffset}`,
    );
    if (!response.ok) throw new Error('データの取得に失敗しました');

    const data = await response.json();
    totalCount = data.count;

    const tweets = data.tweets;
    if (!tweets) return;

    const tweetList = document.getElementById('tweet-list');
    tweetList.innerHTML = '';

    if (data.tweets) {
      data.tweets.forEach((tweet) => {
        tweetList.appendChild(createTweetCard(tweet));
      });
    }

    console.log(`デバック用: ${data}`);

    updateUI();
  } catch (error) {
    console.error('Error', error);
  }
};

// ツイート投稿
const post = () => {
  document.getElementById('tweet-form').addEventListener('submit', (e) => {
    e.preventDefault();

    const textValue = document.getElementById('tweet-content').value;

    fetch('/post', {
      method: 'post',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ content: textValue }),
    });
  });
};

// ツイート詳細
const getTweet = async () => {
  // 現在のパス (/user-detail/1) から最後の部分を取り出す
  const pathParts = window.location.pathname.split('/');
  const tweetId = pathParts[pathParts.length - 1];
  // let id = urlParams.get('id') || 1;
  const response = await fetch(`/api/tweets/${tweetId}`);
  const data = await response.json();

  container = document.getElementById('tweet-detail-container');
  container.innerHTML = '';
  container.appendChild(createTweetCard(data));
};

const dispatchPathTask = async () => {
  const path = window.location.pathname;
  const pathParts = path.split('/');
  const idFromPath = pathParts[pathParts.length - 1];

  if (path.includes('home')) {
    currentApiUrl = '/api/tweets';
    loadTweets();
    setupPagination();
  } else if (path.includes('user-detail')) {
    // クエリパラメータではなく、パスから取ったIDをチェック
    if (!idFromPath || isNaN(idFromPath)) {
      console.error('IDがパスに含まれていません');
      return;
    }
    // 先にユーザー情報を取得して画面に出す（終わるまで次へ行かない
    // 先にDBにアクセスしてデータの取得に失敗したため
    await getUser();
    currentApiUrl = `/api/users/${idFromPath}/tweets`;
    loadTweets();
    setupPagination();
  } else if (path.includes('user-retweet')) {
    currentApiUrl = `/api/users/${idFromPath}/retweets`;
    loadTweets();
    setupPagination();
  } else if (path.includes('post')) {
    post();
  } else if (path.includes('tweet-detail')) {
    getTweet();
  }
};

dispatchPathTask();
