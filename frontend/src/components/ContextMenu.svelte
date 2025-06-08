<script lang="ts">
  import { createEventDispatcher } from "svelte";
  import { portal } from "../lib/portal"; // 新增：導入 portal
  import type {
    ContextMenuItem,
    ContextMenuConfig,
  } from "../types/contextMenu";

  // 組件屬性
  export let visible: boolean = false;
  export let x: number = 0;
  export let y: number = 0;
  export let type: string = ""; // 'row' | 'column' | 'cell'
  export let menuConfig: ContextMenuConfig = {}; // 由父組件傳入的菜單配置
  export let context: any = {}; // 上下文信息，如行索引、列索引等

  // 創建事件分發器
  const dispatch = createEventDispatcher();

  // 默認菜單項目配置（作為 fallback）
  const defaultMenuItems: ContextMenuConfig = {
    row: [
      { id: "insertRowAbove", label: "在上方插入行", icon: "⬆️" },
      { id: "insertRowBelow", label: "在下方插入行", icon: "⬇️" },
      { id: "separator1", type: "separator" },
      { id: "deleteRow", label: "刪除行", icon: "🗑️", danger: true },
    ],
    column: [
      { id: "insertColumnLeft", label: "在左邊插入欄", icon: "⬅️" },
      { id: "insertColumnRight", label: "在右邊插入欄", icon: "➡️" },
      { id: "separator1", type: "separator" },
      { id: "deleteColumn", label: "刪除欄", icon: "🗑️", danger: true },
    ],
    cell: [
      { id: "copy", label: "複製", icon: "📋" },
      { id: "paste", label: "貼上", icon: "📄" },
      { id: "separator1", type: "separator" },
      { id: "clear", label: "清除內容", icon: "🧹" },
    ],
  };

  // 處理菜單項目點擊
  function handleMenuItemClick(action: string) {
    dispatch("action", { action, context });
  }

  // 處理菜單外部點擊
  function handleBackdropClick() {
    dispatch("close");
  }

  // 點擊 ContextMenu 內部時，阻止事件冒泡到 backdrop
  function handleMenuClick(event: MouseEvent) {
    event.stopPropagation();
  }

  // 獲取當前類型的菜單項目，優先使用父組件傳入的配置
  $: currentMenuItems = (menuConfig[type] ||
    defaultMenuItems[type] ||
    []) as ContextMenuItem[]; // 計算選單位置，避免超出螢幕邊界
  $: menuLeft = (() => {
    if (typeof window === "undefined") return x;
    const menuWidth = 200; // 選單最小寬度
    const padding = 8; // 距離螢幕邊緣的間距

    // 調試信息
    console.log("Menu positioning:", {
      x,
      y,
      windowWidth: window.innerWidth,
      windowHeight: window.innerHeight,
      calculated: { left: x, top: y },
    });

    // 檢查右邊是否會超出螢幕
    if (x + menuWidth + padding > window.innerWidth) {
      // 超出右邊界，顯示在滑鼠左側
      return Math.max(padding, x - menuWidth);
    } // 正常顯示在滑鼠位置（不偏移）
    return Math.max(padding, x);
  })();
  $: menuTop = (() => {
    if (typeof window === "undefined") return y;
    // 計算實際選單高度，基於實際的 CSS 變數值
    const itemHeight = 40; // padding: 8px top + 8px bottom + ~24px 內容高度
    const separatorHeight = 17; // 1px 分隔線 + 2 * 8px margin
    const menuPadding = 8; // padding: 4px top + 4px bottom (--spacing-extra-small)
    const screenPadding = 8; // 距離螢幕邊緣的間距

    let menuHeight = menuPadding;
    currentMenuItems.forEach((item) => {
      menuHeight += item.type === "separator" ? separatorHeight : itemHeight;
    });

    // 檢查下方是否會超出螢幕
    if (y + menuHeight + screenPadding > window.innerHeight) {
      // 超出下邊界，顯示在滑鼠上方
      return Math.max(screenPadding, y - menuHeight);
    } // 正常顯示在滑鼠位置（不偏移）
    return Math.max(screenPadding, y);
  })();
</script>

{#if visible}
  <!-- 透明背景，用於捕獲外部點擊 -->
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <!-- svelte-ignore a11y-no-static-element-interactions -->
  <div class="context-menu-backdrop" on:click={handleBackdropClick}>
    <!-- 右鍵菜單 -->
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-static-element-interactions -->
    <div
      use:portal
      class="context-menu"
      style="left: {menuLeft}px; top: {menuTop}px;"
      on:click={handleMenuClick}
    >
      {#each currentMenuItems as item}
        {#if item.type === "separator"}
          <div class="menu-separator"></div>
        {:else}
          <!-- svelte-ignore a11y-click-events-have-key-events -->
          <!-- svelte-ignore a11y-no-static-element-interactions -->
          <div
            class="menu-item"
            class:danger={item.danger}
            class:disabled={item.disabled}
            on:click={() => !item.disabled && handleMenuItemClick(item.id)}
          >
            <span class="menu-icon">{item.icon || ""}</span>
            <span class="menu-label">{item.label || ""}</span>
          </div>
        {/if}
      {/each}
    </div>
  </div>
{/if}

<style>
  .context-menu-backdrop {
    position: fixed;
    top: 0;
    left: 0;
    width: 100vw;
    height: 100vh;
    z-index: 1000;
    background: transparent;
  }

  .context-menu {
    position: absolute;
    background: var(--surface-glass);
    backdrop-filter: blur(var(--blur-heavy));
    border: 1px solid var(--border-color);
    border-radius: var(--radius-medium);
    box-shadow: var(--shadow-4);
    padding: var(--spacing-extra-small) 0;
    min-width: 200px;
    font-family: var(--font-primary);
    font-size: 14px;
    z-index: 1001;
    animation: contextMenuSlideIn var(--transition-duration) ease-out;
    overflow: hidden;
  }

  .context-menu::before {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 1px;
    background: linear-gradient(
      90deg,
      transparent,
      rgba(255, 255, 255, 0.3),
      transparent
    );
  }

  .menu-item {
    display: flex;
    align-items: center;
    padding: var(--spacing-small) var(--spacing-medium);
    cursor: pointer;
    transition: all var(--transition-duration);
    position: relative;
    color: var(--text-primary);
    margin: 0 var(--spacing-extra-small);
    border-radius: var(--radius-small);
  }

  .menu-item::before {
    content: "";
    position: absolute;
    left: 0;
    top: 0;
    bottom: 0;
    width: 0;
    background: linear-gradient(135deg, var(--primary-500), var(--primary-600));
    transition: width var(--transition-duration);
    border-radius: var(--radius-small);
  }

  .menu-item:hover {
    background: linear-gradient(135deg, var(--primary-50), var(--primary-100));
    color: var(--primary-700);
    transform: translateX(2px);
    box-shadow: var(--shadow-1);
  }

  .menu-item:hover::before {
    width: 3px;
  }

  .menu-item.danger {
    color: var(--error-600);
  }

  .menu-item.danger::before {
    background: linear-gradient(135deg, var(--error-500), var(--error-600));
  }

  .menu-item.danger:hover {
    background: linear-gradient(135deg, var(--error-50), var(--error-100));
    color: var(--error-700);
  }

  .menu-item.disabled {
    color: var(--text-disabled);
    cursor: not-allowed;
    opacity: 0.5;
  }

  .menu-item.disabled:hover {
    background: transparent;
    transform: none;
    box-shadow: none;
  }

  .menu-item.disabled:hover::before {
    width: 0;
  }

  .menu-icon {
    margin-right: var(--spacing-small);
    font-size: 16px;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 24px;
    height: 24px;
    border-radius: var(--radius-small);
    background: rgba(255, 255, 255, 0.1);
    backdrop-filter: blur(var(--blur-light));
    transition: all var(--transition-duration);
  }

  .menu-item:hover .menu-icon {
    background: rgba(255, 255, 255, 0.2);
    transform: scale(1.1);
  }

  .menu-item.danger:hover .menu-icon {
    background: rgba(244, 67, 54, 0.1);
  }

  .menu-label {
    flex: 1;
    font-weight: 500;
  }

  .menu-separator {
    height: 1px;
    background: linear-gradient(
      90deg,
      transparent,
      var(--border-color),
      transparent
    );
    margin: var(--spacing-extra-small) var(--spacing-medium);
    position: relative;
  }

  .menu-separator::after {
    content: "";
    position: absolute;
    top: 1px;
    left: 0;
    right: 0;
    height: 1px;
    background: linear-gradient(
      90deg,
      transparent,
      rgba(255, 255, 255, 0.1),
      transparent
    );
  }

  /* 動畫 */
  @keyframes contextMenuSlideIn {
    from {
      opacity: 0;
      transform: scale(0.95) translateY(-10px);
    }
    to {
      opacity: 1;
      transform: scale(1) translateY(0);
    }
  }

  /* 響應式設計 */
  @media (max-width: 480px) {
    .context-menu {
      min-width: 160px;
      font-size: 13px;
    }

    .menu-item {
      padding: var(--spacing-extra-small) var(--spacing-small);
    }

    .menu-icon {
      width: 20px;
      height: 20px;
      font-size: 14px;
    }
  }
</style>
