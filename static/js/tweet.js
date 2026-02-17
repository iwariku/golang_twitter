const urlParams = new URLSearchParams(window.location.search);
const LIMIT = 10;
let currentOffset = parseInt(urlParams.get('offset')) || 0;
let totalCount = 0;

async function loadTweets(offset = 0) {
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
}

function updateUI() {
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
}

document.addEventListener('DOMContentLoaded', () => {
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
});
