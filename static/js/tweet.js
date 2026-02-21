const urlParams = new URLSearchParams(window.location.search);
const LIMIT = 10;
let currentOffset = parseInt(urlParams.get('offset')) || 0;
let totalCount = 0;

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

const loadTweetList = () => {
  const loadTweets = async (offset = 0) => {
    try {
      currentOffset = offset;

      const newUrl = `${window.location.pathname}?LIMIT=${LIMIT}&offset=${currentOffset}`;
      window.history.pushState({ offset: currentOffset }, '', newUrl);

      const response = await fetch(
        `/api/tweets?LIMIT=${LIMIT}&offset=${currentOffset}`,
      );
      if (!response.ok) throw new Error('データの取得に失敗しました');

      const data = await response.json();
      totalCount = data.count;

      const tweets = data.tweets;
      if (!tweets) return;

      const tweetList = document.getElementById('tweet-list');
      tweetList.innerHTML = '';
      tweets.forEach((tweet) => {
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
        </div>
        </div>`;
        tweetList.appendChild(tweetCard);
      });

      updateUI();
    } catch (error) {
      console.error('Error', error);
    }
  };

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

  loadTweets(currentOffset);
};

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

const getTweet = async () => {
  const urlParams = new URLSearchParams(window.location.search); // 共通にする
  let currentParams = urlParams.get('id') || 1; // 共通にする
  const response = await fetch(`/api/tweet-detail?id=${currentParams}`);

  // 画面表示の部分
  const data = await response.json();
  console.log('サーバーから届いたデータ:', data);
  container = document.getElementById('tweet-detail-container');
  container.innerHTML = `
    <div class="w-[600px] min-w-[600px] border-x border-gray-200 min-h-screen">
      <div class="p-4 border-b border-gray-200">
        <div class="flex items-center mb-4">
          <div class="w-12 h-12 bg-gray-200 rounded-full mr-3"></div>
          <div>
            <div class="font-bold text-[15px]">ユーザーID: ${data.user_id}</div>
          </div>
        </div>

        <div class="text-[23px] leading-8 whitespace-pre-wrap break-words mb-4">
          ${data.content}
        </div>
      </div>
    </div>
  `;
};

const dispatchPathTask = () => {
  const path = window.location.pathname;
  if (path.includes('home')) {
    loadTweetList();
  } else if (path.includes('post')) {
    post();
  } else if (path.includes('detail')) {
    getTweet();
  }
};

dispatchPathTask();
