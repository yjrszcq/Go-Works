$(document).ready(function() {
    let currentUser = null;
    let currentRole = null;
    showMenuPage();
    hideAllNav();
    showNav("menu-nav");
    showNav("login-nav");


    // 导航栏点击事件
    $('#menu-nav').click(function() {
        showMenuPage();
    });

    $('#cart-nav').click(function() {
        if (currentRole) {
            showCartPage();
        } else {
            showLogin();
        }
    });

    $('#employee-nav').click(function() {
        showEmployeePage();
    });

    $('#deliveryman-nav').click(function() {
        showDeliveryPage();
    });

    $('#login-nav').click(function() {
        showLoginPage();
        //test();
        //showMyOtherSection();
        //showMyPage();
       // showAdminPage();
        //showSignupPage();
        //showEmployeePage();
        //showLoginPage();
        //showDeliveryPage()
    });

    $('#my-nav').click(function() {
        $("#profile-nav").click();
        showMyPage();
    });


    $('#admin-nav').click(function() {
        if ( currentRole === 'admin') {
            showAdminPage();
        } else {
            showLogin();
        }
    });


   
    // 登录
    $('#login-form').submit(function(event) {
        event.preventDefault();
        login();
    });


    function showNav(para_nav){
        $("#"+para_nav).show();
    }

    function hideNav(para_nav){
        $("#"+para_nav).hide();
    }

    function hideAllNav(){
        hideNav("menu-nav");
        hideNav("cart-nav");
        hideNav("employee-nav");
        hideNav("deliveryman-nav");
        hideNav("my-nav");
        hideNav("login-nav");
        hideNav("admin-nav");

    }

    function showMyPage(){
        hideAllPages();
        $('#my-section').show();
    }

    function showMenuPage() {
        hideAllPages();
        $('#menu-section').show();
        loadMenuItems();
    }

    function showCartPage() {
        hideAllPages();
        $('#cart-section').show();
        getCartList()
       
    }

    function showChangePasswordPage(){
        hideAllPages();
        $('#change-password-section').show();
    }

    function showAdminPage() {
        hideAllPages();
        $('#admin-section').show();
    }
    
    function showSignupPage(){
        hideAllPages();
        $('#signup-section').show();
    }

    function showLoginPage(){
        hideAllPages();
        $('#login-section').show();
    }

    function showEmployeePage(){
        hideAllPages();
        $('#employee-section').show();
    }


    function showDeliveryPage(){
        hideAllPages();
        $('#deliveryman-section').show();
        loadOrders()
    }

    function hideAllPages() {
        $('.page').hide();
        $('#error-message').hide();
        $('#success-message').hide();
    }

    function loadMenuItems() {
        $.ajax({
            url: 'http://127.0.0.1:1000/dish/list',
            method: 'GET',
            xhrFields: {
                withCredentials: true // 带上 cookie
            },
            success: function(data) {
                renderMenuItems(data);
            },
            error: function() {
                showError('加载菜单失败');
            }
        });
    }
    
    function searchMenuItems() {
        let searchName = $('#dish-name').val().trim();
        if (searchName) {
            $.ajax({
                url: 'http://127.0.0.1:1000/dish/find/name',
                method: 'POST',
                contentType: 'application/json',
                data: JSON.stringify({ name: searchName }),
                xhrFields: {
                    withCredentials: true // 带上 cookie
                },
                success: function(data) {
                    renderMenuItems(data);
                },
                error: function() {
                    showError('搜索菜品失败');
                }
            });
        } else {
            showError('请输入菜品名字');
        }
    }
    
    function renderMenuItems(data) {
        let menuHtml = '';
        data.forEach(function(dish) {
            menuHtml += `
                <div class="dish" data-id="${dish.Id}">
                    <img src="${dish.ImageURL}" alt="${dish.Name}">
                    <div class="dish-info">
                        <p>ID：${dish.Id}</p>
                        <p>菜名：${dish.Name}</p>
                        <p>价格：${dish.Price}</p>
                    </div>
                    <div class="dish-actions">
                        <button class="view-details">查看详情</button>
                    </div>
                </div>
            `;
        });
        $('#menu-items').html(menuHtml);
        
    }

    $('#search-button').click(searchMenuItems);

    // 绑定输入框回车事件
    $('#dish-name').keypress(function(e) {
        if (e.which == 13) { // 13 是回车键的键码
            searchMenuItems();
        }
    });
    

     
    // 如果需要绑定查看详情按钮的点击事件，也可以使用事件委托
    $('#menu-items').on('click', '.view-details', function() {
        let dishId = $(this).closest('.dish').data('id');
        if (currentRole) {
            showDishDetails(dishId);
        } else {
            showLoginPage();
        }
    });
        
    function showDishDetails(dishId) {
        hideAllPages();
        $('#dishdetail-section').show();
        
        // 发送 AJAX 请求获取菜品详细信息
        $.ajax({
            url: 'http://127.0.0.1:1000/dish/find/id', // 请替换成你的实际 API 地址
            method: 'POST',
            data: JSON.stringify({id: dishId }),
            dataType: 'json',
            xhrFields: {
                withCredentials: true
            },
            success: function(response) {
                // 将返回的数据填充到相应的元素中
                $('#dish-image').attr('src', response.ImageURL);
                $('#dish-detail-id').text(response.Id);
                $('#dish-detail-name').text(response.Name);
                $('#dish-detail-description').text(response.Description);
                $('#dish-detail-price').text('价格: ' + response.Price + '元');

                // 发送 AJAX 请求获取评论
                getCategoryById(response.CategoryID);
                getDishReviews(dishId);
            },
            error: function(error) {
                console.log('请求失败', error);
                alert('获取菜品详细信息失败，请稍后再试。');
            }
         });
   
    function getCategoryById(categoryId){
        $.ajax({
            url: 'http://127.0.0.1:1000/category/find/id',
            method: 'POST',
            xhrFields: {
                withCredentials: true
            },
            contentType: 'application/json',
            data: JSON.stringify({id: categoryId}),
            dataType: 'json',
            success: function(data) {
                $("#dish-detail-category").text('种类为: ' + data.Name);
            },
            error: function() {
                $("#dish-detail-category").text('未获取到');
            }
        });
    }

    function getDishReviews(dishId){
        $.ajax({
            url: 'http://127.0.0.1:1000/review/find/dish_id',
            method: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({ dish_id: dishId }),
            success: function(response) {
                renderReviews(response);
            },
            error: function(error) {
                console.error('请求失败:', error);
            }
        });
    }

    
    function renderReviews(reviews) {
        const $reviewsContainer = $('#dish-reviews');
        $reviewsContainer.empty(); // 清空现有评论

        reviews.forEach(function(review) {
            const reviewHtml = `
                <div class="review">
                    <div class="review-header">
                        <p class="review-id">评论 ID: ${review.Id}</p>
                        <div class="review-user">
                            <p class="review-customer">用户 ID: ${review.CustomerID}</p>
                            <p class="review-rating">${review.Rating} 星</p>
                        </div>
                    </div>
                    <p class="review-comment">${review.Comment}</p>
                    <p class="review-date">评论日期: ${review.ReviewDate}</p>
                </div>
            `;
            $reviewsContainer.append(reviewHtml);
        });
    }


      

        $('#back-to-menu').click(function() {
            // 返回上一级餐单的逻辑
            hideAllPages();
            $('#menu-section').show(); 
        });
    
    }
    
    $('#add-to-cart').off('click').on('click', function(e) {
        e.preventDefault(); // 阻止表单的默认提交行为
        // 加入购物车的逻辑
        addToCart();
    });

    function addToCart() {
        let dishId = $('#dish-detail-id').text();
        let quantity = $('#quantity').val();

        $.ajax({
            url: 'http://127.0.0.1:1000/cart/add',
            type: 'POST',
            xhrFields: {
                withCredentials: true
            },
            contentType: 'application/json',
            data: JSON.stringify({dish_id:Number(dishId), quantity:Number(quantity) }),
            dataType: 'json',
            success: function() {
                alert('已将 ' + $('#dish-name').text() + ' 加入购物车。');
            },
            error: function(jqXHR, textStatus, errorThrown) {
                console.error('加入购物车失败:', textStatus, errorThrown);
                alert('加入购物车失败');
            }
        });
    }

    function getCartList() {

        $.ajax({
            url: 'http://127.0.0.1:1000/cart/list',
            type: 'GET',
            xhrFields: {
                withCredentials: true
            },
            dataType: 'json',
            success: function(data) {
                let cartItems = $('#cart-items');
                cartItems.empty();
                data.forEach(cartItem => {
                    cartItems.append(`
                        <div class="cart-item">
                            <div class="cart-item-info">
                                <p>订单ID: ${cartItem.Id}</p>
                                <p>客户ID: ${cartItem.CustomerID}</p>
                                <p>菜品ID: ${cartItem.DishID}</p>
                                <p>数量: ${cartItem.Quantity}</p>
                                <p>加入时间: ${cartItem.CreatedAt}</p>
                            </div>
                            <div class="cart-item-actions">
                                <button class="delete-cart-item" data-cart-id="${cartItem.Id}">删除</button>
                            </div>
                        </div>
                    `);
                });

                // 为“删除”按钮添加点击事件
                $('.delete-cart-item').on('click', function() {
                    e.preventDefault(); // 阻止表单的默认提交行为
                    deleteCartItem($(this).data('cart-id'));
                    
                });
            },
            error: function(jqXHR, textStatus, errorThrown) {
                console.error('获取购物车失败:', textStatus, errorThrown);
            }
        });
    }

    function deleteCartItem(cartItemId) {
        $.ajax({
            url: 'http://127.0.0.1:1000/cart/delete',
            type: 'POST',
            xhrFields: {
                withCredentials: true
            },
            contentType: 'application/json',
            data: JSON.stringify({ id: cartItemId }),
            dataType: 'json',
            success: function(data) {
                alert('成功删除购物车项');
                getCartList();
            },
            error: function(jqXHR, textStatus, errorThrown) {
                console.error('删除购物车项失败:', textStatus, errorThrown);
                alert('删除购物车项失败');
            }
        });
    }

    function clearCart() {
        $.ajax({
            url: 'http://127.0.0.1:1000/cart/clear',
            type: 'GET',
            xhrFields: {
                withCredentials: true
            },
            dataType: 'json',
            success: function(data) {
                alert('成功清除全部购物车');
                getCartList();
            },
            error: function(jqXHR, textStatus, errorThrown) {
                console.error('清除全部购物车失败:', textStatus, errorThrown);
                alert('清除全部购物车失败');
            }
        });
    }


    // 为“清除全部购物车”按钮添加点击事件
    $('#clear-cart').on('click', function() {
        clearCart();
    });

    // 为“确认支付”按钮添加点击事件
    $('#confirm-payment').on('click', function() {
        // 未来可以在这里添加支付逻辑
        let cartItems = [];
        $('#cart-items .cart-item').each(function() {
            let cartItemId = $(this).find('.delete-cart-item').data('cart-id');
            cartItems.push(cartItemId);
        });

        let deliveryLocation = $('#delivery-location').val();
        let deliveryTime = $('#delivery-time').val() .replace('T', ' ')+":00";

        let data = {
            delivery_location: deliveryLocation,
            delivery_time: deliveryTime,
            cart_item_id: cartItems
        };

        $.ajax({
            url: 'http://127.0.0.1:1000/order/create',
            type: 'POST',
            contentType: 'application/json',
            data: JSON.stringify(data),
            xhrFields: {
                withCredentials: true
            },
            success: function(response) {
                let orderId = response.Id;
                loadPaymentDetails(orderId);
            },
            error: function(error) {
                alert('创建订单失败: ' + error.responseText);
            }
        });

    });

     // 加载支付详情页面
     function loadPaymentDetails(orderId) {
        $('#cart-section').hide();
        $('#payment-section').show();
        $('#order-items').data('order-id', orderId); // 存储订单ID
        $.ajax({
            url: 'http://127.0.0.1:1000/order_item/find/order_id',
            type: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({ order_id: orderId }),
            xhrFields: {
                withCredentials: true
            },
            success: function(response) {
                let total = 0;
                $('#order-items').empty();
                response.forEach(item => {
                    total += item.TotalPrice;
                    $('#order-items').append(`
                        <div class="order-item">
                            <p>订单项ID: ${item.Id}</p>
                            <p>订单ID: ${item.OrderID}</p>
                            <p>菜品ID: ${item.DishID}</p>
                            <p>数量: ${item.Quantity}</p>
                            <p>单价: ${item.UnitPrice} 元</p>
                            <p>总金额: ${item.TotalPrice} 元</p>
                            <p>创建时间: ${item.CreatedAt}</p>
                        </div>
                    `);
                });
                $('#total-price').text(total.toFixed(2));
            },
            error: function(error) {
                alert('获取订单详情失败: ' + error.responseText);
            }
        });
    }

    // 支付订单按钮点击事件
    $('#pay-order').on('click', function() {
        let orderId = $('#order-items').data('order-id');
        if (orderId) {
            $.ajax({
                url: 'http://127.0.0.1:1000/order/pay',
                type: 'POST',
                contentType: 'application/json',
                data: JSON.stringify({ id: orderId }),
                xhrFields: {
                    withCredentials: true
                },
                success: function(response) {
                    alert('支付成功');
                    showMenuPage();
                },
                error: function(error) {
                    alert('支付失败: ' + error.responseText);
                    showMenuPage();
                }
            });
        } else {
            alert('订单ID未找到');
        }
    });

    // 取消订单按钮点击事件
    $('#cancel-order').on('click', function() {
        let orderId = $('#order-items').data('order-id');
        if (orderId) {
            $.ajax({
                url: 'http://127.0.0.1:1000/order/cancel',
                type: 'POST',
                contentType: 'application/json',
                data: JSON.stringify({ id: orderId }),
                xhrFields: {
                    withCredentials: true
                },
                success: function(response) {
                    alert('取消成功');
                    showHomePage()
                },
                error: function(error) {
                    alert('取消失败: ' + error.responseText);
                    showHomePage()
                }
            });
        } else {
            alert('订单ID未找到');
        }
    });

    // 关闭页面按钮点击事件
    $('#close-payment').on('click', function() {
        showMenuPage();
    });


    

   
    /*
    function showError(message) {
        $('#error-message p').text('错误信息：' + message);
        $('#error-message').show();
    }
    
    function showSuccess(message) {
        $('#success-message p').text('成功信息：' + message);
        $('#success-message').show();
    }
    */


   
    $('#toggle-password').click(function() {
        if ($('#my-section #password').attr('type') === 'password') {
            $('#my-section #password').attr('type', 'text');
            $(this).text('隐藏密码');
        } else {
            $('#my-section #password').attr('type', 'password');
            $(this).text('显示密码');
        }
    });

    // 确认修改功能
    $('#submit-profile').click(function(e) {
        e.preventDefault(); // 阻止表单的默认提交行为

        // 获取表单数据
        const name = $('#profile-name').val();
        const phone = $('#profile-phone').val();
        const email = $('#profile-email').val();
        const address = $('#profile-address').val(); 

        // 验证所有字段都不为空
        if (!name || !phone || !email || !address) {
            alert('所有字段都必须填写！');
            return;
        }
        // 根据全局变量 currentRole 决定调用哪个API端点
        let url;
        if (currentRole === 'employee') {
            url = 'http://127.0.0.1:1000/employee/edit';
        } else if (currentRole === 'customer') {
            url = 'http://127.0.0.1:1000/customer/edit';
        } else {
            alert('未知角色，无法提交！');
            return;
        }
    
        // 构建要发送的数据
        const data = {
            name: name,
            phone: phone,
            email: email,
            address : address
        };
    
        // 使用 jQuery 的 AJAX 方法提交数据
        $.ajax({
            url: url,
            xhrFields: {
                withCredentials: true
            },
            type: 'POST',
            data: JSON.stringify(data),
            contentType: 'application/json',
            dataType: 'json', // 设置 dataType 为 json
            success: function(response) {
                alert('提交成功！');
                console.log(response);
            },
            error: function(jqXHR, textStatus, errorThrown) {
                try {
                    var errorResponse = JSON.parse(jqXHR.responseText);
                    alert(errorResponse.message);
                    console.error(errorResponse);
                } catch (e) {
                    alert('请求失败: ' + textStatus + ', ' + errorThrown);
                    console.error(jqXHR);
                }
            }
        });
    
        // 显示表单数据（仅用于调试）
       // alert('姓名: ' + name + '\n电话: ' + phone + '\n电子邮件: ' + email + '\n密码: ' + password);
    });







  
    $('.nav-link').click(function(e) {
        e.preventDefault();
        $('.nav-link').removeClass('active');
        $(this).addClass('active');
    
        const target = $(this).data('target');
    
        $('.content-section').removeClass('active');
        $('#' + target).addClass('active');
    
        // 隐藏所有内容部分
        $('.content-section').hide();
    
        // 显示当前激活的内容部分
        $('#' + target).show();
    
        if (target === 'orders') {
            getOrders();
        } else if (target === 'reviews') {
            getMyReviews();
        }
    });

    function getOrders(){
        $.ajax({
            url: 'http://127.0.0.1:1000/order/list',
            type: 'GET',
            xhrFields: {
                withCredentials: true // 带上 cookie
            },
            dataType: 'json',
            success: function(data) {
                let ordersTable = $('#orders-table tbody');
                ordersTable.empty(); // 清空现有的订单数据
                data.forEach(order => {
                    ordersTable.append(`
                        <tr>
                            <td>${order.Id}</td>
                            <td>${order.OrderDate}</td>
                            <td>${order.Status}</td>
                            <td>
                                <button type="button" class="order-details" data-order-id="${order.Id}">查看详情</button>
                            </td>
                        </tr>
                    `);
                });

                // 为“查看详情”按钮添加点击事件
                $('.order-details').on('click', function() {
                    let orderId = $(this).data('order-id');
                    getMyOrderById(orderId);
                });
            },
            error: function(jqXHR, textStatus, errorThrown) {
                console.error('请求失败:', textStatus, errorThrown);
            }
        });
    }

    function getMyOrderById(orderId) {
        $.ajax({
            url: 'http://127.0.0.1:1000/order/find/id',
            type: 'POST', // 修改为 POST 请求
            xhrFields: {
                withCredentials: true // 带上 cookie
            },
            contentType: 'application/json',
            data: JSON.stringify({ id: orderId }),
            dataType: 'json',
            success: function(data) {
                showOrderDetails(data);
            },
            error: function(jqXHR, textStatus, errorThrown) {
                console.error('请求失败:', textStatus, errorThrown);
            }
        });
    }
    
    function showOrderDetails(order) {
        $('#order-details').addClass('active');
        $('#details-order-items').empty();
         // 隐藏所有内容部分
         $('.content-section').hide();
         // 显示当前激活的内容部分
         $('#order-details').show();
        // 创建新界面
        let orderDetailsSection = `
                    <div class="form-group">
                        <strong>订单ID:</strong> ${order.Id}<br>
                        <strong>用户ID:</strong> ${order.CustomerID}<br>
                        <strong>配送地址:</strong> ${order.DeliveryLocation}<br>
                        <strong>订单状态:</strong> ${order.Status}<br>
                        <strong>支付状态:</strong> ${order.PaymentStatus}<br>
                        <strong>订单总额:</strong> ${order.TotalAmount}<br>
                        <strong>订单日期:</strong> ${order.OrderDate}<br>
                        <strong>配送时间:</strong> ${order.DeliveryTime}<br>
                        <strong>配送人员ID:</strong> ${order.DeliveryPersonID}<br>
                        <strong>创建时间:</strong> ${order.CreatedAt}<br>
                         <strong>如果订单已送达，即可评价:</strong><br>
                        <button type="button" id="back-to-orders" class="button-container">返回订单列表</button>
                        ${order.Status === '待支付' ? '<button type="button" id="confirm-payment" class="button-container">确认支付</button>' : ''}
                        
                    </div>      
        `
       
        fetchOrderItems(order.Id,order.Status);
        $('#details-order-items').append(orderDetailsSection);
    
        // 为“返回订单列表”按钮添加点击事件
        $('#back-to-orders').on('click', function() {
            $('.content-section').hide();
            // 显示当前激活的内容部分
            $('#orders-nav').click();
           // 创建新界面
        });
    }

    $('#edit-password').on('click', function() {
        showChangePasswordSection();
    });

    function showChangePasswordSection() {
        hideAllPages();
        showChangePasswordPage();
    }

    

    // 绑定返回登录按钮的点击事件
    $('#c-back-to-login').click(function() {
       showMyPage();
    });

    // 绑定修改密码表单的提交事件
    $('#change-password-form').submit(function(event) {
        event.preventDefault(); // 阻止表单默认提交行为

        const oldPassword = $('#old-password').val();
        const newPassword = $('#new-password').val();
        const confirmPassword = $('#c-confirm-password').val();

        // 检查新密码和确认密码是否一致
        if (newPassword !== confirmPassword) {
            alert('新密码和确认密码不一致，请重新输入。');
            return;
        }

        // 构建请求数据
        const data = {
            old_password: oldPassword,
            new_password: newPassword,
            confirm_password: confirmPassword
        };

        // 根据当前角色发送请求
        const url = currentRole === 'customer' ? 'http://127.0.0.1:1000/customer/change_password' : 'http://127.0.0.1:1000/employee/change_password';

        $.ajax({
            url: url,
            method: 'POST',
            contentType: 'application/json',
            data: JSON.stringify(data),
            xhrFields: {
                withCredentials: true // 带上 cookie
            },
            success: function(response) {
                console.log('Password changed successfully:', response);
                alert('密码修改成功！');
                // 你可以在这里显示登录界面或进行其他操作
                showMyPage();
            },
            error: function(jqXHR, textStatus, errorThrown) {
                try {
                    var errorResponse = JSON.parse(jqXHR.responseText);
                    alert(errorResponse.message);
                    console.error(errorResponse);
                } catch (e) {
                    alert('请求失败: ' + textStatus + ', ' + errorThrown);
                    console.error(jqXHR);
                }
            }
        });
    });


    function fetchOrderItems(orderId, orderStatus) {
        $.ajax({
            url: 'http://127.0.0.1:1000/order_item/find/order_id',
            method: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({ order_id: orderId }),
            dataType: 'json',
            xhrFields: {
                withCredentials: true // 带上 cookie
            },
            success: function(response) {
                if (response && response.length > 0) {
                    const orderItemsContainer = $('#order-items-container');
                    orderItemsContainer.empty(); // 清空之前的订单项信息
                    orderItemsContainer.append(`
                        <p>当前订单的Id为: ${orderId}</p>
                    `);
                    response.forEach(o => {
                        orderItemsContainer.append(`
                            <div class="order-item" data-order-item-id="${o.Id}">
                                <p>当前订单的订单项Id为: ${o.Id}</p>
                                <p>当前订单的餐品为: ${o.DishName}</p>
                                <p>当前订单的状态为: ${o.ReviewStatus}</p>
                                ${orderStatus === "已送达" ? '<button class="add-review">添加评论</button>' : ''}
                            </div>
                        `);
                    });
    
                    // 为每个订单项的添加评论按钮绑定点击事件
                    $('.order-item .add-review').off('click').on('click', function() {
                        const orderItemId = $(this).closest('.order-item').data('order-item-id');
                        $('#order-item-id').val(orderItemId);
                        $('#order-id').val(orderId);
    
                        // 在提交评论按钮点击时获取当前输入的值
                        $('#submit-review').off('click').on('click', function() {
                            const reviewRating = Number($("#review-rating").val()); // 获取评分的值
                            const reviewComment = $("#review-comment").val(); // 获取评论的值
                            createReview(orderItemId, reviewRating, reviewComment);
                        });
    
                        $('#review-form').show();
                    });
                } else {
                    $('#order-items-container').html('<p>没有找到订单项信息。</p>');
                }
            },
            error: function(jqXHR, textStatus, errorThrown) {
                console.error('请求失败:', textStatus, errorThrown);
                $('#order-items-container').html('<p>请求失败，请稍后再试。</p>');
            }
        });
    }
    
    function createReview(reviewId, reviewRating, reviewComment) {
        $.ajax({
            url: 'http://127.0.0.1:1000/review/create',
            method: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({
                order_item_id: reviewId,
                rating: reviewRating,
                comment: reviewComment
            }),
            xhrFields: {
                withCredentials: true // 带上 cookie
            },
            success: function(response) {
                console.log('Review created successfully:', response);
                // 可以在这里隐藏评论表单或显示成功提示
                $('#review-form').hide();
                alert('评论已成功提交！');
            },
            error: function(error) {
                console.error('Error creating review:', error);
                alert('评论提交失败，请稍后再试。');
            }
        });
    }
   

    function getMyReviews() {
        $.ajax({
            url: 'http://127.0.0.1:1000/review/list',
            type: 'GET',
            xhrFields: {
                withCredentials: true // 带上 cookie
            },
            dataType: 'json',
            success: function(data) {
                let reviewsList = $('#reviews-list');
                reviewsList.empty(); // 清空现有的评论数据
                data.forEach(review => {
                    reviewsList.append(`
                        <li>
                            <div class="review-item">
                                <strong>评价ID:</strong> ${review.Id}<br>
                                <strong>客户ID:</strong> ${review.CustomerID}<br>
                                <strong>菜品ID:</strong> ${review.DishID}<br>
                                <strong>评分:</strong> ${review.Rating}<br>
                                <strong>评价内容:</strong> ${review.Comment}<br>
                                <strong>评价日期:</strong> ${review.CreatedAt}<br>
                            </div>
                        </li>
                    `);
                });
            },
            error: function(jqXHR, textStatus, errorThrown) {
                console.error('请求失败:', textStatus, errorThrown);
            }
        });
    }
    
   
    

    $('input[name="role"]').change(function() {
        let currentRole = $(this).val();
        if (currentRole === 'admin') {
            $('#email-form-group').hide();
            $('#username-form-group').show();
            $('#login-email').attr('name', '');
            $('#login-username').attr('name', 'username');
            $('#login-email').val('');
            $('#login-username').val('');
        } else {
            $('#email-form-group').show();
            $('#username-form-group').hide();
            $('#login-email').attr('name', 'email');
            $('#login-username').attr('name', '');
            $('#login-email').val('');
            $('#login-username').val('');
        }
    });


    

    //登录有关的操作
    function login() {
        let email = $('#login-email').val();
        let username = $('#login-username').val();
        let password = $('#login-password').val();
        currentRole = $('input[name="role"]:checked').val();

        let data = {};
        if (currentRole === 'admin') {
            data = { name: username, password: password };
        } else {
            data = { email: email, password: password };
        }

        let url = currentRole === 'employee' 
            ? '/employee/login' 
            : currentRole === 'admin' 
            ? '/admin/login' 
            : '/customer/login';

        $.ajax({
            url: 'http://127.0.0.1:1000' + url,
            method: 'POST',
            xhrFields: {
                withCredentials: true
            },
            contentType: 'application/json',
            data: JSON.stringify(data),
            success: function() {
                $('#login-email').val('');
                $('#login-password').val('');
                $('#login-username').val('');
                $('input[name="role"]').prop('checked', false);
                if (currentRole === 'admin') {
                    hideAllPages();
                    hideAllNav();
                    showNav("admin-nav");
                    showAdminPage();
                } else {
                    getProfile(currentRole);
                }
            },
            error: function() {
                $('#login-email').val('');
                $('#login-password').val('');
                $('#login-username').val('');
                $('input[name="role"]').prop('checked', false);
                showError('登录失败');
                currentRole = null;
            }
        });
    }
    
    function getProfile(role) {
        $.ajax({
            url: 'http://127.0.0.1:1000/'+role+'/profile',
            method: 'GET',
            xhrFields: {
                withCredentials: true
            },
        success: function(data) {
            // 设置表单字段的值
            $('#profile-name').val(data.Name);
            $('#profile-phone').val(data.Phone);
            $('#profile-email').val(data.Email);
            $('#profile-address').val(data.Address); 

            // 根据角色跳转到不同的页面
            if(role==="customer"){
                hideAllNav();
                showNav("menu-nav");
                showNav("cart-nav");
                showNav("my-nav");
                showOtherSection();
                showMenuPage();
            }else{
                 hideOtherSection();               
                if (data.Role === '员工') {
                    hideAllNav();
                    showNav("employee-nav")
                    showNav("my-nav");
                    showEmployeePage();
                } else if(data.Role === '送餐员') {
                    hideAllNav();
                    showNav("my-nav");
                    showNav("deliveryman-nav");
                    showDeliveryPage() ;
                } 
            }
            
        },
        error: function() {
            alert('获取用户信息失败');
        }
        });
    }

   function hideOtherSection(){
        $("#orders-nav").hide();
        $("#reviews-nav").hide();
   }

   function showOtherSection(){
        $("#orders-nav").show();
        $("#reviews-nav").show();
    }

    // 取消登录功能
    $('#logout').click(function() {
            logout();
    });

    function logout(){
        if (!confirm("确定要登出吗？")) {
            return;
        }

        const logoutUrl = currentRole === 'employee' 
        ? '/employee/logout' 
        : currentRole === 'admin' 
        ? '/admin/logout' 
        : '/customer/logout';

        $.ajax({
            url: 'http://127.0.0.1:1000'+logoutUrl,
            method: 'GET',
            xhrFields: {
                withCredentials: true // 确保带上 cookies
            },
            success: function(response) {
                // 处理服务器返回的数据
                console.log('Logout successful:', response);
                showSuccess('已注销');
    
                // 清空全局变量
                currentRole = null;
    
                // 显示登录导航，隐藏个人导航
                hideAllNav();
                showNav('login-nav');
                showNav('menu-nav')
                hideNav('my-nav');
    
                // 重定向到登录页面
                showLoginPage();
            },
            error: function(error) {
                console.error('Error logging out:', error);
                showSuccess('注销失败，请重试');
            }
        });
    }


    //登陆成功提示
    function showSuccess(message) {
        alert(message);
    }

    // 登录失败提示
    function showError(message) {
        alert(message);
    }

    

  
   
    $('#submit-login').click(function() {
        login();
    });

    $('#register').click(function() {
        showSignupPage();
    });

    $('#back-to-login').click(function() {
        showLoginPage();
    });

   
    $('#employee-fields').hide();

    // 监听角色选择变化
    $('input[name="signup-role"]').change(function() {
        const role = $(this).val();
        if (role === 'employee') {
            $('#employee-fields').show();
        } else {
            $('#employee-fields').hide();
            // 清空员工字段
            $('#signup-name').val('');
            $('#signup-phone').val('');
        }
    });

     // 监听提交按钮点击
     $('#submit-signup').click(function(event) {
        event.preventDefault();

        const email = $('#signup-email').val();
        const password = $('#signup-password').val();
        const confirmPassword = $('#signup-confirm-password').val();
        const role = $('input[name="signup-role"]:checked').val();
        const name = role === 'employee' ? $('#signup-name').val() : '';
        const phone = role === 'employee' ? $('#signup-phone').val() : '';

        if (password !== confirmPassword) {
            alert('密码和确认密码不匹配');
            return;
        }

        const url = role === 'employee' ? '/employee/signup' : '/customer/signup';
        const data = {
            email: email,
            password: password,
            confirm_password: confirmPassword
        };

        if (role === 'employee') {
            data.name = name;
            data.phone = phone;
        }

        $.ajax({
            url: 'http://127.0.0.1:1000'+url,
            type: 'POST',
            contentType: 'application/json',
            data: JSON.stringify(data),
            xhrFields: {
                withCredentials: true // 确保带上 cookies
            },
            success: function(response) {
                $('#signup-email').val('');
                $('#signup-password').val('');
                $('#signup-confirm-password').val('');
                $('input[name="signup-role"]').prop('checked', false);
                if (role === 'employee') {
                    $('#signup-name').val('');
                    $('#signup-phone').val('');
                }
                alert('注册成功');
                showLoginPage();
            },
            error: function(jqXHR, textStatus, errorThrown) {
                try {
                    var errorResponse = JSON.parse(jqXHR.responseText);
                    alert(errorResponse.message);
                    console.error(errorResponse);
                } catch (e) {
                    alert('请求失败: ' + textStatus + ', ' + errorThrown);
                    console.error(jqXHR);
                    $('#signup-email').val('');
                    $('#signup-password').val('');
                    $('#signup-confirm-password').val('');
                    $('input[name="signup-role"]').prop('checked', false);
                    if (role === 'employee') {
                        $('#signup-name').val('');
                        $('#signup-phone').val('');
                    }
                }
            }
    
        });
    });

    // 监听返回登录按钮点击
    $('#back-to-login').click(function() {
        showLoginPage();
    });

    //管理员部分
    // 获取所有员工信息
    function fetchAllEmployees() {
        $.ajax({
            url: 'http://127.0.0.1:1000/admin/employee/list',
            method: 'GET',
            xhrFields: {
                withCredentials: true // 确保带上 cookies
            },
            success: function(data) {
                if (data && typeof data === 'object' && !Array.isArray(data)) {
                    data = [data];
                }
                displayEmployees(data);
            },
            error: function(xhr, status, error) {
                console.error('获取员工信息失败:', error);
                alert('获取员工信息失败: ' + error);
            }
        });
    }

    // 根据状态获取员工信息
    function fetchEmployeesByStatus(status) {
        $.ajax({
            url: 'http://127.0.0.1:1000/admin/employee/find/status',
            method: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({ status: status }),
            xhrFields: {
                withCredentials: true // 确保带上 cookies
            },
            success: function(data) {
                if (data && typeof data === 'object' && !Array.isArray(data)) {
                    data = [data];
                }
                displayEmployees(data);
            },
            error: function(xhr, status, error) {
                console.error('获取员工信息失败:', error);
                alert('获取员工信息失败: ' + error);
            }
        });
    }

    // 根据ID查询员工信息
    function fetchEmployeeById(id) {
        $.ajax({
            url: 'http://127.0.0.1:1000/admin/employee/find/id',
            method: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({ id: id }),
            xhrFields: {
                withCredentials: true // 确保带上 cookies
            },
            success: function(data) {
                if (data && typeof data === 'object' && !Array.isArray(data)) {
                    data = [data];
                }
                displayEmployees(data);
            },
            error: function(xhr, status, error) {
                console.error('获取员工信息失败:', error);
                alert('获取员工信息失败: ' + error);
            }
        });
    }

    // 根据姓名查询员工信息
    function fetchEmployeeByName(name) {
        $.ajax({
            url: 'http://127.0.0.1:1000/admin/employee/find/name',
            method: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({ name: name }),
            xhrFields: {
                withCredentials: true // 确保带上 cookies
            },
            success: function(data) {
                if (data && typeof data === 'object' && !Array.isArray(data)) {
                    data = [data];
                }
                displayEmployees(data);
            },
            error: function(xhr, status, error) {
                console.error('获取员工信息失败:', error);
                alert('获取员工信息失败: ' + error);
            }
        });
    }

    // 修改员工角色
    function editEmployeeRole(id, role) {
        $.ajax({
            url: 'http://127.0.0.1:1000/admin/employee/edit/role',
            method: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({ id: id, role: role }),
            xhrFields: {
                withCredentials: true // 确保带上 cookies
            },
            success: function(data) {
                alert('修改角色成功');
                fetchAllEmployees();
            },
            error: function(xhr, status, error) {
                console.error('修改角色失败:', error);
                alert('修改角色失败: ' + error);
            }
        });
    }

    // 修改员工状态
    function editEmployeeStatus(id, status) {
        $.ajax({
            url: 'http://127.0.0.1:1000/admin/employee/edit/status',
            method: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({ id: id, status: status }),
            xhrFields: {
                withCredentials: true // 确保带上 cookies
            },
            success: function(data) {
                alert('修改状态成功');
                fetchAllEmployees();
            },
            error: function(xhr, status, error) {
                console.error('修改状态失败:', error);
                alert('修改状态失败: ' + error);
            }
        });
    }

    // 显示员工信息
    function displayEmployees(employees) {
        const employeeList = $('.employee-list');
        employeeList.empty();

        employees.forEach(employee => {
            let employeeItem = `
                <div class="employee-item">
                    <strong>员工ID:</strong> ${employee.Id}<br>
                    <strong>姓名:</strong> ${employee.Name}<br>
                    <strong>角色:</strong> ${employee.Role}<br>
                    <strong>邮箱:</strong> ${employee.Email}<br>
                    <strong>电话:</strong> ${employee.Phone}<br>
                    <strong>状态:</strong> ${employee.Status}<br>
                    <strong>创建时间:</strong> ${employee.CreatedAt}<br>
                    <strong>更新时间:</strong> ${employee.UpdatedAt}<br>
                </div>
            `;

            employeeList.append(employeeItem);

           
        });
    }

    $('#all-employees').on('click', function() {
        fetchAllEmployees();
    });

    // 为“可用”按钮添加点击事件
    $('#available-employees').on('click', function() {
        fetchEmployeesByStatus('可用');
    });

    // 为“不可用”按钮添加点击事件
    $('#unavailable-employees').on('click', function() {
        fetchEmployeesByStatus('不可用');
    });

    // 为“ID查询”按钮添加点击事件
    $('#search-by-id').on('click', function() {
        const id = parseInt($('#id-search').val());
        fetchEmployeeById(id);
    });

    // 为“姓名查询”按钮添加点击事件
    $('#search-by-name').on('click', function() {
        const name = $('#name-search').val();
        fetchEmployeeByName(name);
    });

    // 为“选择角色”按钮添加点击事件
    $('#edit-role').on('click', function() {
        const id = parseInt($('#id-edit-role').val());
        const role = $('#role-select').val();
        editEmployeeRole(id, role);
    });

    // 为“选择状态”按钮添加点击事件
    $('#edit-status').on('click', function() {
        const id = parseInt($('#id-edit-status').val());
        const status = $('#status-select').val();
        editEmployeeStatus(id, status);
    });

    $('#fetch-customers').click(function(e) {
        e.preventDefault();
        $.ajax({
            url: 'http://127.0.0.1:1000/admin/customer/list',
            type: 'GET',
            xhrFields: {
                withCredentials: true // 确保带上 cookies
            },
            success: function(response) {
                if (response && typeof response === 'object' && !Array.isArray(response)) {
                    response = [response];
                }
                displayCustomers(response);
            },
            error: function(error) {
                alert('获取顾客信息失败，请重试！');
                console.error(error);
            }
        });
    });

    // 根据顾客ID查询顾客信息
    $('#find-customer-by-id').click(function(e) {
        e.preventDefault();
        const customerId = parseInt($('#id-find-customer').val());
        if (!customerId) {
            alert('请输入顾客ID！');
            return;
        }
        $.ajax({
            url: 'http://127.0.0.1:1000/admin/customer/find/id',
            type: 'POST',
            data: JSON.stringify({id: customerId}),
            contentType: 'application/json',
            xhrFields: {
                withCredentials: true // 确保带上 cookies
            },
            success: function(response) {
                if (response && typeof response === 'object' && !Array.isArray(response)) {
                    response = [response];
                }
                displayCustomers(response);
            },
            error: function(error) {
                alert('查询顾客信息失败，请重试！');
                console.error(error);
            }
        });
    });

    // 根据顾客姓名查询顾客信息
    $('#find-customer-by-name').click(function(e) {
        e.preventDefault();
        const customerName = $('#name-find-customer').val();
        if (!customerName) {
            alert('请输入顾客姓名！');
            return;
        }
        $.ajax({
            url: 'http://127.0.0.1:1000/admin/customer/find/name',
            type: 'POST',
            data: JSON.stringify({name: customerName}),
            contentType: 'application/json',
            xhrFields: {
                withCredentials: true // 确保带上 cookies
            },
            success: function(response) {
                displayCustomers(response);
            },
            error: function(error) {
                alert('查询顾客信息失败，请重试！');
                console.error(error);
            }
        });
    });

    // 重置顾客密码
    $('#reset-password').click(function(e) {
        e.preventDefault();
        const customerId = parseInt($('#id-reset-password').val());
        const new_password = $('#password-reset').val();
        if (!customerId || !new_password) {
            alert('请输入顾客ID和新密码！');
            return;
        }
        $.ajax({
            url: 'http://127.0.0.1:1000/admin/customer/edit/init_password',
            type: 'POST',
            data: JSON.stringify({id: customerId, password: new_password}),
            contentType: 'application/json',
            xhrFields: {
                withCredentials: true // 确保带上 cookies
            },
            success: function(response) {
                alert('密码重置成功！');
                console.log(response);
            },
            error: function(jqXHR, textStatus, errorThrown) {
                try {
                    var errorResponse = JSON.parse(jqXHR.responseText);
                    alert(errorResponse.message);
                    console.error(errorResponse);
                } catch (e) {
                    alert('请求失败: ' + textStatus + ', ' + errorThrown);
                    console.error(jqXHR);
                }
            }
        });
    });

    // 显示所有顾客信息
    function displayCustomers(customers) {
        const customerList = $('.customer-list');
        customerList.empty(); // 清空之前的顾客信息
        if (customers.length === 0) {
            customerList.append('<p>暂无顾客信息</p>');
            return;
        }
        customers.forEach(function(customer) {
            customerList.append(`
                <div class="customer-item">
                    <strong>顾客ID:</strong> ${customer.Id}<br>
                    <strong>姓名:</strong> ${customer.Name}<br>
                    <strong>邮箱:</strong> ${customer.Email}<br>
                    <strong>电话:</strong> ${customer.Phone}<br>
                    <strong>地址:</strong> ${customer.Address}<br>
                    <strong>创建时间:</strong> ${customer.CreatedAt}<br>
                    <strong>更新时间:</strong> ${customer.UpdatedAt}<br>
                </div>
            `);
        });
    }

    // 显示单个顾客信息
    function displayCustomer(customer) {
        const customerList = $('.customer-list');
        customerList.empty(); // 清空之前的顾客信息
        if (!customer) {
            customerList.append('<p>未找到顾客信息</p>');
            return;
        }
        customerList.append(`
            <div class="customer-item">
                <strong>顾客ID:</strong> ${customer.CustomerID}<br>
                <strong>姓名:</strong> ${customer.Name}<br>
                <strong>邮箱:</strong> ${customer.Email}<br>
                <strong>电话:</strong> ${customer.Phone}<br>
                <strong>地址:</strong> ${customer.Address}<br>
                <strong>创建时间:</strong> ${customer.CreatedAt}<br>
                <strong>更新时间:</strong> ${customer.UpdatedAt}<br>
            </div>
        `);
    }

    $('#admin-logout').click(function() {
        logout()
    });


    //员工界面
    //employee-section操作
    $('#show-all-dishes').click(function() {
        $.ajax({
            url: 'http://127.0.0.1:1000/dish/list',
            method: 'GET',
            xhrFields: {
                withCredentials: true // 带上 cookie
            },
            success: function(response) {
                const dishesList = $('#dishes-list');
                dishesList.empty();

                response.forEach(function(dish) {
                    dishesList.append(`
                        <li class="dish-item">
                            <img src="${dish.ImageURL}" alt="${dish.Name}">
                            <div class="details">
                                <strong>${dish.Name}</strong>
                                <p>菜品ID：${dish.Id}<p>
                                <p>价格: ${dish.Price}</p>
                                <p>描述: ${dish.Description}</p>
                                <p>类别ID: ${dish.CategoryID}</p>
                            </div>
                        </li>
                    `);
                });
            },
            error: function() {
                showError('加载菜单失败');
            }
        });   
    });

    $('#find-by-id').click(function() {
        const id = Number($('#find-id').val());
        fetchDishes('/dish/find/id', { id: id });
    });

    $('#find-by-name').click(function() {
        const name = $('#find-name').val();
        fetchDishes('/dish/find/name', { name: name });
    });

    $('#find-by-category').click(function() {
        const category = Number($('#find-category').val());
        fetchDishes('/dish/find/category', { category_id: category});
    });

    $('#add-dish').click(function() {
        const name = $('#add-name').val();
        const imageUrl = $('#add-image-url').val();
        const price = parseFloat($('#add-price').val());
        const description = $('#add-description').val();
        const categoryId = Number($('#add-category-id').val());

        $.ajax({
            url: 'http://127.0.0.1:1000/employee/dish/create',
            type: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({
                name: name,
                image_url: imageUrl,
                price: price,
                description: description,
                category_id: categoryId
            }),
            xhrFields: {
                withCredentials: true // 带上 cookie
            },
            success: function(response) {
                alert('菜品创建成功');
                
            },
            error: function(xhr, status, error) {
                alert('菜品创建失败: ' + error);
            }
        });
    });

    $('#edit-dish').click(function() {
        const id = parseInt($('#edit-id').val());
        const name = $('#edit-name').val();
        const imageUrl = $('#edit-url').val();
        const price = parseFloat($('#edit-price').val());
        const description = $('#edit-description').val();
        const categoryId = parseInt($('#edit-category').val());


        $.ajax({
            url: 'http://127.0.0.1:1000/employee/dish/edit',
            type: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({
                id:id,
                name: name,
                image_url: imageUrl,
                price: price,
                description: description,
                category_id: categoryId
            }),
            xhrFields: {
                withCredentials: true // 带上 cookie
            },
            success: function(response) {
                alert('菜品修改成功');
            },
            error: function(xhr, status, error) {
                alert('菜品修改失败: ' + error);
            }
        });
    });

    function fetchDishes(url, data = {}) {
        $.ajax({
            url: 'http://127.0.0.1:1000' + url,
            type: 'POST',
            contentType: 'application/json',
            data: JSON.stringify(data),
            xhrFields: {
                withCredentials: true // 带上 cookie
            },
            success: function(response) {
                console.log('Response:', response); // 输出响应数据
                const dishesList = $('#dishes-list');
                if (dishesList.length === 0) {
                    alert('未找到 #dishes-list 元素');
                    return;
                }
                dishesList.empty();

                // 确保 response 是一个数组
                const dishes = Array.isArray(response) ? response : [response];

                dishes.forEach(function(dish) {
                    dishesList.append(`
                        <li class="dish-item">
                            <img src="${dish.ImageURL}" alt="${dish.Name}">
                            <div class="details">
                                <strong>${dish.Name}</strong>
                                <p>价格: ${dish.Price}</p>
                                <p>描述: ${dish.Description}</p>
                                <p>类别ID: ${dish.CategoryID}</p>
                            </div>
                        </li>
                    `);
                });
            },
            error: function(xhr, status, error) {
                alert('获取菜品列表失败: ' + error);
            }
        });
    }


    $('#show-all-orders').click(function() {
        $.ajax({
            url: 'http://127.0.0.1:1000/employee/order/list',
            method: 'GET',
            xhrFields: {
                withCredentials: true // 带上 cookie
            },
            success: function(data) {
                updateOrderList(data);
            },
            error: function() {
                $('#order-list').html('<li>获取订单列表失败</li>');
            }
        });
    });

    // 按配送员ID查找订单
    $('#find-delivery-person-id').click(function() {
        var deliverymanId = parseInt($('#find-by-delivery-person-id').val());
        if (deliverymanId) {
            $.ajax({
                url: 'http://127.0.0.1:1000/employee/order/find/delivery_person_id',
                method: 'POST',
                data: JSON.stringify({ deliveryman_id: deliverymanId }),
                contentType: 'application/json',
                xhrFields: {
                    withCredentials: true // 带上 cookie
                },
                success: function(data) {
                    updateOrderList(data);
                },
                error: function() {
                    $('#order-list').html('<li>查找订单失败</li>');
                }
            });
        } else {
            alert('请输入配送员ID');
        }
    });

    // 按客户ID查找订单
    $('#find-customer-id').click(function() {
        var customerId = parseInt($('#find-by-customer-id').val());
        if (customerId) {
            $.ajax({
                url: 'http://127.0.0.1:1000/employee/order/find/customer_id',
                method: 'POST',
                data: JSON.stringify({ customer_id: customerId }),
                contentType: 'application/json',
                xhrFields: {
                    withCredentials: true // 带上 cookie
                },
                success: function(data) {
                    updateOrderList(data);
                },
                error: function() {
                    $('#order-list').html('<li>查找订单失败</li>');
                }
            });
        } else {
            alert('请输入客户ID');
        }
    });

    // 取消订单
    $('#cancel-order-d').click(function() {
        var orderId =parseInt($('#cancel-order-id').val());
        if (orderId) {
            $.ajax({
                url: 'http://127.0.0.1:1000/employee/order/cancel',
                method: 'POST',
                data: JSON.stringify({ id: orderId }),
                contentType: 'application/json',
                xhrFields: {
                    withCredentials: true // 带上 cookie
                },
                success: function(data) {
                    alert('订单取消成功');
                    $('#order-list').html('');
                    // 可以重新获取订单列表
                    $('#show-all-orders').click();
                },
                error: function(jqXHR, textStatus, errorThrown) {
                    try {
                        var errorResponse = JSON.parse(jqXHR.responseText);
                        alert(errorResponse.message);
                        console.error(errorResponse);
                    } catch (e) {
                        alert('请求失败: ' + textStatus + ', ' + errorThrown);
                        console.error(jqXHR);
                    }
                }
            });
        } else {
            alert('请输入订单ID');
        }
    });

    // 确认订单
    $('#confirm-order-d').click(function() {
        var orderId = parseInt($('#confirm-order-id').val());
        if (orderId) {
            $.ajax({
                url: 'http://127.0.0.1:1000/employee/order/confirm',
                method: 'POST',
                data: JSON.stringify({ id: orderId }),
                contentType: 'application/json',
                xhrFields: {
                    withCredentials: true // 带上 cookie
                },
                success: function(data) {
                    alert('订单确认成功');
                    $('#order-list').html('');
                    // 可以重新获取订单列表
                    $('#show-all-orders').click();
                },
                error: function(jqXHR, textStatus, errorThrown) {
                    try {
                        var errorResponse = JSON.parse(jqXHR.responseText);
                        alert(errorResponse.message);
                        console.error(errorResponse);
                    } catch (e) {
                        alert('请求失败: ' + textStatus + ', ' + errorThrown);
                        console.error(jqXHR);
                    }
                }
            });
        } else {
            alert('请输入订单ID');
        }
    });

    // 完成订单
    $('#complete-order-d').click(function() {
        var orderId = parseInt($('#complete-order-id').val());
        if (orderId) {
            $.ajax({
                url: 'http://127.0.0.1:1000/employee/order/complete',
                method: 'POST',
                data: JSON.stringify({ id: orderId }),
                contentType: 'application/json',
                xhrFields: {
                    withCredentials: true // 带上 cookie
                },
                success: function(data) {
                    alert('订单完成成功');
                    $('#order-list').html('');
                    // 可以重新获取订单列表
                    $('#show-all-orders').click();
                },
                error: function(jqXHR, textStatus, errorThrown) {
                    try {
                        var errorResponse = JSON.parse(jqXHR.responseText);
                        alert(errorResponse.message);
                        console.error(errorResponse);
                    } catch (e) {
                        alert('请求失败: ' + textStatus + ', ' + errorThrown);
                        console.error(jqXHR);
                    }
                }
            });
        } else {
            alert('请输入订单ID');
        }
    });

    // 更新订单列表
    function updateOrderList(orders) {
        var html = '';
        orders.forEach(function(order) {
             // 检查 data 是否为对象，如果是，则将其转换为数组
            html += `
                <li class="order-item">
                    <div class="details">
                        <p>ID: ${order.Id}</p>
                        <p>客户ID: ${order.CustomerID}</p>
                        <p>配送地点: ${order.DeliveryLocation}</p>
                        <p>状态: ${order.Status}</p>
                        <p>支付状态: ${order.PaymentStatus}</p>
                        <p>总金额: ${order.TotalAmount}</p>
                        <p>订单日期: ${order.OrderDate}</p>
                        <p>配送时间: ${order.DeliveryTime}</p>
                        <p>配送员ID: ${order.DeliveryPersonID}</p>
                        <p>创建时间: ${order.CreatedAt}</p>
                        <p>更新时间: ${order.UpdatedAt}</p>
                    </div>
                </li>
            `;
        });
        $('#order-list').html(html);
    }

    


    
    function loadOrders() {
        $.ajax({
            url: 'http://127.0.0.1:1000/deliveryman/order/find/waiting_for_delivery',
            method: 'GET',
            xhrFields: {
                withCredentials: true // 带上 cookie
            },
            success: function(data) {
                // 检查 data 是否为对象，且包含 message 属性
                if (data && typeof data === 'object' && data.message) {
                    if (data.message === '成功, 暂无待送订单') {
                        renderOrders([]);
                        return;
                    }
                    // 如果有其他成功消息，可以根据需要处理
                    alert(data.message);
                    return;
                }
    
                // 检查 data 是否为对象，如果是，则将其转换为数组
                if (data && typeof data === 'object' && !Array.isArray(data)) {
                    data = [data];
                }
    
                renderOrders(data);
            },
            error: function() {
                showError('加载订单失败');
            }
        });
    }
    
    function renderOrders(orders) {
        $('#orders-list').empty();
        let ordersHtml = '';
    
        if (orders.length === 0) {
            // 没有订单的情况
            ordersHtml = `<p>当前没有订单。</p>`;
        } else {
            // 多个订单（包括单个订单转换后的数组）的情况
            orders.forEach(function(order) {
                ordersHtml += `
                    <div class="order" data-id="${order.Id}">
                        <div class="order-info">
                            <p>订单ID：${order.Id}</p>
                            <p>客户ID：${order.CustomerID}</p>
                            <p>配送地址：${order.DeliveryLocation}</p>
                            <p>状态：${order.Status}</p>
                            <p>下单日期：${order.OrderDate}</p>
                            <p>预计配送时间：${order.DeliveryTime}</p>
                        </div>
                        <div class="order-actions">
                            <button class="deliver-order">出餐</button>
                            <button class="delivered-order">已送达</button>
                        </div>
                    </div>
                `;
            });
        }
    
        $('#orders-list').html(ordersHtml);
    
        // 绑定出餐按钮点击事件
        $('.deliver-order').click(function() {
            let orderId = parseInt($(this).closest('.order').data('id'));
            updateOrderStatus('http://127.0.0.1:1000/deliveryman/order/deliver', orderId, '出餐成功');
        });
    
        // 绑定已送达按钮点击事件
        $('.delivered-order').click(function() {
            let orderId = parseInt($(this).closest('.order').data('id'));
            updateOrderStatus('http://127.0.0.1:1000/deliveryman/order/delivered', orderId, '已送达成功');
        });
    }
    
    
    function updateOrderStatus(url, orderId, successMessage) {
        $.ajax({
            url: url,
            method: 'POST',
            contentType: 'application/json',
            data: JSON.stringify({ id: orderId }),
            xhrFields: {
                withCredentials: true // 带上 cookie
            },
            success: function() {
                showError(successMessage);
                loadOrders(); // 重新加载订单列表
            },
            error: function(jqXHR, textStatus, errorThrown) {
                try {
                    var errorResponse = JSON.parse(jqXHR.responseText);
                    alert(errorResponse.message);
                    console.error(errorResponse);
                } catch (e) {
                    alert('请求失败: ' + textStatus + ', ' + errorThrown);
                    console.error(jqXHR);
                }
            }
        });
    }

});

