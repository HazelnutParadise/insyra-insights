<script lang="ts">
  import { createEventDispatcher } from "svelte";
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

  // 獲取當前類型的菜單項目，優先使用父組件傳入的配置
  $: currentMenuItems = (menuConfig[type] ||
    defaultMenuItems[type] ||
    []) as ContextMenuItem[];
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
      class="context-menu"
      style="left: {x}px; top: {y}px;"
      on:click|stopPropagation
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
    background: white;
    border: 1px solid #ddd;
    border-radius: 6px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    padding: 4px 0;
    min-width: 180px;
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Oxygen,
      Ubuntu, Cantarell, "Open Sans", "Helvetica Neue", sans-serif;
    font-size: 14px;
    z-index: 1001;
  }
  .menu-item {
    display: flex;
    align-items: center;
    padding: 8px 12px;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .menu-item:hover {
    background-color: #f5f5f5;
  }

  .menu-item.danger {
    color: #d32f2f;
  }

  .menu-item.danger:hover {
    background-color: #ffebee;
  }

  .menu-item.disabled {
    color: #bbb;
    cursor: not-allowed;
  }

  .menu-item.disabled:hover {
    background-color: transparent;
  }

  .menu-icon {
    margin-right: 8px;
    font-size: 16px;
    display: flex;
    align-items: center;
    justify-content: center;
    width: 20px;
  }

  .menu-label {
    flex: 1;
  }

  .menu-separator {
    height: 1px;
    background-color: #e0e0e0;
    margin: 4px 0;
  }
</style>
